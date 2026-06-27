package alerting

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"

	"github.com/tfdriftctl/tfdriftctl/internal/config"
	"github.com/tfdriftctl/tfdriftctl/internal/model"
)

const emailHeadersTmpl = "Subject: [tfdriftctl Alert] High Risk Drift Detected in %s\r\n" +
	"To: %s\r\n" +
	"From: %s\r\n" +
	"MIME-version: 1.0;\r\n" +
	"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n"

const emailHTML = `<html>
<body>
	<h2>tfdriftctl Alert: High Risk Drift</h2>
	<p>A drift scan has completed with a Total Risk Score of <strong>{{.Report.Summary.TotalRiskScore}}</strong>.</p>
	<ul>
		<li><strong>Workspace:</strong> {{.Report.Workspace}}</li>
		<li><strong>Scan ID:</strong> {{.Report.ScanID}}</li>
		<li><strong>Findings:</strong> {{.Report.Summary.TotalFindings}}</li>
	</ul>
	
	<h3>High Risk Findings</h3>
	<ul>
		{{range .Report.Findings}}
			{{if ge .RiskScore 70}}
				<li><strong>{{.Kind}}</strong> - {{.ResourceType}} ({{.ResourceName}}): Risk {{.RiskScore}}</li>
			{{end}}
		{{end}}
	</ul>
	<p>Please check the tfdriftctl dashboard for full details.</p>
</body>
</html>
`

// SendDriftAlert sends an email if the drift report exceeds the configured minimum risk score.
func SendDriftAlert(report *model.DriftReport, cfg config.AlertingConfig) {
	if !cfg.Enabled || report.Summary.TotalRiskScore < cfg.MinimumRiskScore {
		return
	}

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)

	headers := fmt.Sprintf(emailHeadersTmpl, report.Workspace, cfg.To, cfg.From)

	tmpl, err := template.New("email").Parse(emailHTML)
	if err != nil {
		log.Printf("alerting: failed to parse email template: %v", err)
		return
	}

	data := struct {
		Report    *model.DriftReport
		To        string
		From      string
		Workspace string
	}{
		Report:    report,
		To:        cfg.To,
		From:      cfg.From,
		Workspace: report.Workspace,
	}

	var body bytes.Buffer
	body.WriteString(headers)

	if err := tmpl.Execute(&body, data); err != nil {
		log.Printf("alerting: failed to execute email template: %v", err)
		return
	}

	err = smtp.SendMail(addr, auth, cfg.From, []string{cfg.To}, body.Bytes())
	if err != nil {
		log.Printf("alerting: failed to send email: %v", err)
		return
	}

	log.Printf("alerting: successfully sent drift alert to %s for workspace %s", cfg.To, report.Workspace)
}
