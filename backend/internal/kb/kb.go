package kb

// Article defines the structure for a knowledge base article.
type Article struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// GetArticles returns a hardcoded slice of articles to simulate a real knowledge base.
func GetArticles() []Article {
	return []Article{
		{
			ID:      "kb-001",
			Title:   "How to reset your password",
			Content: "To reset your password, go to the login page and click on the 'Forgot Password' link. You will receive an email with instructions on how to set a new password. Make sure to choose a strong password that you haven't used before.",
		},
		{
			ID:      "kb-002",
			Title:   "VPN Connection Issues",
			Content: "If you are having trouble connecting to the company VPN, first ensure you have the latest version of the VPN client installed. Second, check your internet connection to make sure it is stable. If the problem persists, try restarting your computer. Contact IT support if you are still unable to connect.",
		},
		{
			ID:      "kb-003",
			Title:   "Setting up a new printer",
			Content: "To set up a new printer, first connect it to the network via an ethernet cable or Wi-Fi. Then, go to your computer's system settings, find the 'Printers & Scanners' section, and click 'Add Printer'. Your computer should automatically detect the printer. If not, you may need to install drivers from the manufacturer's website.",
		},
	}
}
