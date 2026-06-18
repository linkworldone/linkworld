package services

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/wneessen/go-mail"
)

// Mailer 发送事务性邮件（当前仅验证码）。从 env 读 SMTP 配置：
//
//	SMTP_HOST / SMTP_PORT / SMTP_USERNAME / SMTP_PASSWORD / SMTP_FROM
//
// 降级：SMTP_HOST 为空时不真发，改 log.Printf("[DEV] ...") 醒目打印验证码，
// 本地不配 SMTP 也能跑完绑定流程。配了 SMTP_HOST 即走真实发送（STARTTLS + PLAIN）。
//
// 安全：SMTP_PASSWORD 绝不进日志；日志只打 host/port/from/to（密码留在内存，仅作 SMTP 鉴权用）。
type Mailer struct {
	host     string
	port     int
	username string
	password string
	from     string
}

// NewMailer 从环境变量装配 Mailer。host 留空即 [DEV] 日志降级（不报错）。
// port 解析失败回退 587（STARTTLS 标准端口）。from 缺省回退 username。
func NewMailer() *Mailer {
	port := 587
	if p := os.Getenv("SMTP_PORT"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			port = v
		}
	}
	from := os.Getenv("SMTP_FROM")
	username := os.Getenv("SMTP_USERNAME")
	if from == "" {
		from = username
	}
	m := &Mailer{
		host:     os.Getenv("SMTP_HOST"),
		port:     port,
		username: username,
		password: os.Getenv("SMTP_PASSWORD"),
		from:     from,
	}
	if m.host == "" {
		log.Println("INFO: 未配置 SMTP_HOST，邮件走 [DEV] 日志降级（验证码打印到日志，不真发）")
	} else {
		// 注意：绝不打印 password
		log.Printf("INFO: Mailer 已装配 host=%s port=%d from=%s", m.host, m.port, m.from)
	}
	return m
}

// SendVerificationCode 给 to 发送 6 位验证码邮件（纯文本，中英双语，10 分钟有效）。
// 降级：host 为空时仅 log.Printf("[DEV] ...") 打印 to+code，返回 nil（不报错）。
func (m *Mailer) SendVerificationCode(to, code string) error {
	subject := "LinkWorld 邮箱验证码 / Verification Code"
	body := fmt.Sprintf(
		"【LinkWorld】您的邮箱验证码是：%s，10 分钟内有效，请勿泄露给他人。\n\n"+
			"[LinkWorld] Your verification code is: %s. It is valid for 10 minutes. Do not share it with anyone.\n",
		code, code,
	)

	// 降级：未配 SMTP_HOST → 只打日志，不真发（密码绝不进日志）。
	if m.host == "" {
		log.Printf("[DEV] 邮箱验证码 to=%s code=%s", to, code)
		return nil
	}

	msg := mail.NewMsg()
	if err := msg.From(m.from); err != nil {
		return fmt.Errorf("set from: %w", err)
	}
	if err := msg.To(to); err != nil {
		return fmt.Errorf("set to: %w", err)
	}
	msg.Subject(subject)
	msg.SetBodyString(mail.TypeTextPlain, body)

	opts := []mail.Option{
		mail.WithPort(m.port),
		mail.WithTLSPolicy(mail.TLSMandatory),
	}
	if m.username != "" {
		opts = append(opts,
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername(m.username),
			mail.WithPassword(m.password),
		)
	}

	client, err := mail.NewClient(m.host, opts...)
	if err != nil {
		return fmt.Errorf("new mail client: %w", err)
	}
	if err := client.DialAndSend(msg); err != nil {
		// 不打 password；err 由 go-mail 产生，不含明文密码。
		log.Printf("ERROR: 发送验证码邮件失败 host=%s from=%s to=%s err=%v", m.host, m.from, to, err)
		return fmt.Errorf("send verification email: %w", err)
	}
	log.Printf("INFO: 验证码邮件已发送 host=%s from=%s to=%s", m.host, m.from, to)
	return nil
}
