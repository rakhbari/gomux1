package utils

type Version struct {
    Timestamp string `json:"timestamp"`
    GitSha    string `json:"gitSha"`
    GitBranch string `json:"gitBranch"`
}
