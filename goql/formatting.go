package goql

import "strings"

// normalizeSQL removes extra spaces and line breaks for comparison
func normalizeSQL(s string) string {
	return strings.Join(strings.Fields(s), " ") // Split into words & rejoin with single spaces
}

// formatSQL attempts to format an SQL string by adding line breaks and spacing
func formatSQL(query string) string {
	// Define common SQL keywords for formatting
	keywords := []string{
		"SELECT", "FROM", "WHERE", "ORDER BY", "GROUP BY", "HAVING", "INSERT INTO",
		"VALUES", "UPDATE", "SET", "DELETE FROM", "JOIN", "LEFT JOIN", "RIGHT JOIN",
		"INNER JOIN", "OUTER JOIN", "LIMIT",
	}

	// Normalize spaces
	query = strings.TrimSpace(query)
	query = strings.Join(strings.Fields(query), " ") // Removes extra spaces

	// Convert SQL keywords to uppercase and insert newlines
	for _, kw := range keywords {
		upperKW := strings.ToUpper(kw)
		query = strings.ReplaceAll(query, upperKW, "\n"+upperKW)
	}

	// Ensure no leading/trailing whitespace per line
	lines := strings.Split(query, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}

	// Join with proper indentation
	return strings.Join(lines, "\n")
}
