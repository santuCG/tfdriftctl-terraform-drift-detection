package alerting

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"text/template"

	"github.com/tfdriftctl/tfdriftctl/internal/config"
	"github.com/tfdriftctl/tfdriftctl/internal/model"
)

const emailTemplate = `Subject: [tfdriftctl Alert] High Risk Drift Detected in {{.Workspace}}
To: {{.To}}
From: {{.From}}
MIME-version: 1.0;
Content-Type: text/html; charset="UTF-8";

<html>
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

	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		log.Printf("alerting: failed to parse email template: %v", err)
		return
	}

	data := struct {
		Report *model.DriftReport
		To     string
		From   string
		Workspace string
	}{
		Report:    report,
		To:        cfg.To,
		From:      cfg.From,
		Workspace: report.Workspace,
	}

	var body bytes.Buffer
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
