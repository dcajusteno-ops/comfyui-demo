package backend

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func normalizePromptAssistantState(state PromptAssistantState) PromptAssistantState {
	state.FavoriteIDs = limitUniqueNonEmptyStrings(state.FavoriteIDs, 512)
	state.RecentIDs = limitUniqueNonEmptyStrings(state.RecentIDs, 120)
	state.ActiveSource = strings.TrimSpace(state.ActiveSource)
	state.ActiveCategory = strings.TrimSpace(state.ActiveCategory)
	state.ActiveSubcategory = strings.TrimSpace(state.ActiveSubcategory)
	state.ActiveScope = strings.TrimSpace(state.ActiveScope)

	switch state.ViewMode {
	case "all", "favorites", "recent":
	default:
		state.ViewMode = "all"
	}

	switch state.ActiveEditor {
	case "positive", "negative":
	default:
		state.ActiveEditor = "positive"
	}

	switch state.ItemsPerPage {
	case 8, 12, 24, 48:
	default:
		state.ItemsPerPage = 12
	}

	if state.CurrentPage < 1 {
		state.CurrentPage = 1
	}

	return state
}

func limitUniqueNonEmptyStrings(values []string, limit int) []string {
	if limit <= 0 {
		limit = len(values)
	}
	result := uniqueNonEmptyStrings(values)
	if len(result) > limit {
		result = result[:limit]
	}
	return result
}

func removeStringValue(values []string, target string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if strings.TrimSpace(value) == "" || value == target {
			continue
		}
		result = append(result, value)
	}
	return uniqueNonEmptyStrings(result)
}

func normalizePromptTextKey(value string) string {
	cleaned := strings.ToLower(strings.TrimSpace(value))
	if cleaned == "" {
		return ""
	}
	return strings.Join(strings.Fields(cleaned), " ")
}

func normalizePromptLibraryEntry(entry PromptLibraryEntry) PromptLibraryEntry {
	entry.ID = strings.TrimSpace(entry.ID)
	entry.Source = strings.TrimSpace(entry.Source)
	entry.Category = strings.TrimSpace(entry.Category)
	entry.Subcategory = strings.TrimSpace(entry.Subcategory)
	entry.Scope = strings.TrimSpace(entry.Scope)
	entry.TextEN = strings.TrimSpace(entry.TextEN)
	entry.TextZH = strings.TrimSpace(entry.TextZH)
	entry.Preview = strings.TrimSpace(entry.Preview)
	entry.ExtraID = strings.TrimSpace(entry.ExtraID)
	entry.SearchText = strings.TrimSpace(entry.SearchText)
	if entry.SearchText == "" {
		entry.SearchText = strings.ToLower(strings.TrimSpace(strings.Join([]string{
			entry.Source,
			entry.Category,
			entry.Subcategory,
			entry.Scope,
			entry.TextEN,
			entry.TextZH,
		}, " ")))
	}
	return entry
}

func promptEntriesDuplicate(left, right PromptLibraryEntry) bool {
	leftEN := normalizePromptTextKey(left.TextEN)
	rightEN := normalizePromptTextKey(right.TextEN)
	if leftEN != "" && leftEN == rightEN {
		return true
	}

	leftZH := normalizePromptTextKey(left.TextZH)
	rightZH := normalizePromptTextKey(right.TextZH)
	if leftZH != "" && leftZH == rightZH {
		return true
	}

	return false
}

func (a *App) loadPromptLibrary() ([]PromptLibraryEntry, error) {
	a.promptLibraryMu.RLock()
	if a.promptLibraryLoaded {
		cached := append([]PromptLibraryEntry(nil), a.promptLibraryCache...)
		a.promptLibraryMu.RUnlock()
		return cached, nil
	}
	a.promptLibraryMu.RUnlock()

	a.promptLibraryMu.Lock()
	defer a.promptLibraryMu.Unlock()

	if a.promptLibraryLoaded {
		return append([]PromptLibraryEntry(nil), a.promptLibraryCache...), nil
	}

	data, err := os.ReadFile(a.promptLibraryFile())
	if err != nil {
		return nil, fmt.Errorf("prompt library not found: %w", err)
	}

	var entries []PromptLibraryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("prompt library parse failed: %w", err)
	}

	for index := range entries {
		entries[index] = normalizePromptLibraryEntry(entries[index])
	}

	a.promptLibraryCache = entries
	a.promptLibraryLoaded = true

	return append([]PromptLibraryEntry(nil), entries...), nil
}

func (a *App) GetPromptLibraryEntries() ([]PromptLibraryEntry, error) {
	return a.loadPromptLibrary()
}

func (a *App) loadCustomPromptEntries() ([]PromptLibraryEntry, error) {
	var entries []PromptLibraryEntry
	data, err := os.ReadFile(a.customPromptEntriesFile())
	if err != nil {
		if os.IsNotExist(err) {
			return []PromptLibraryEntry{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []PromptLibraryEntry{}, nil
	}
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	result := make([]PromptLibraryEntry, 0, len(entries))
	for _, entry := range entries {
		entry = normalizePromptLibraryEntry(entry)
		if entry.Source == "" {
			entry.Source = customPromptSource
		}
		if entry.ID == "" {
			entry.ID = uuid.New().String()
		}
		if entry.TextEN == "" && entry.TextZH == "" {
			continue
		}
		result = append(result, entry)
	}
	return result, nil
}

func (a *App) saveCustomPromptEntries(entries []PromptLibraryEntry) error {
	normalized := make([]PromptLibraryEntry, 0, len(entries))
	for _, entry := range entries {
		entry = normalizePromptLibraryEntry(entry)
		if entry.Source == "" {
			entry.Source = customPromptSource
		}
		if entry.ID == "" {
			entry.ID = uuid.New().String()
		}
		if entry.TextEN == "" && entry.TextZH == "" {
			continue
		}
		normalized = append(normalized, entry)
	}
	data, _ := json.MarshalIndent(normalized, "", "  ")
	return os.WriteFile(a.customPromptEntriesFile(), data, 0644)
}

func (a *App) GetCustomPromptEntries() ([]PromptLibraryEntry, error) {
	return a.loadCustomPromptEntries()
}

func (a *App) AddCustomPromptEntry(entry PromptLibraryEntry) (PromptLibraryEntry, error) {
	entry = normalizePromptLibraryEntry(entry)
	entry.Source = customPromptSource
	if entry.TextEN == "" && entry.TextZH == "" {
		return PromptLibraryEntry{}, fmt.Errorf("prompt text is empty")
	}

	systemEntries, err := a.loadPromptLibrary()
	if err != nil {
		return PromptLibraryEntry{}, err
	}
	for _, item := range systemEntries {
		if promptEntriesDuplicate(item, entry) {
			return PromptLibraryEntry{}, fmt.Errorf("系统词库中已存在重复提示词")
		}
	}

	customEntries, err := a.loadCustomPromptEntries()
	if err != nil {
		return PromptLibraryEntry{}, err
	}
	for _, item := range customEntries {
		if promptEntriesDuplicate(item, entry) {
			return PromptLibraryEntry{}, fmt.Errorf("我的词库中已存在重复提示词")
		}
	}

	entry.ID = uuid.New().String()
	customEntries = append([]PromptLibraryEntry{entry}, customEntries...)
	if err := a.saveCustomPromptEntries(customEntries); err != nil {
		return PromptLibraryEntry{}, err
	}
	return entry, nil
}

func (a *App) DeleteCustomPromptEntry(id string) error {
	entries, err := a.loadCustomPromptEntries()
	if err != nil {
		return err
	}

	nextEntries := make([]PromptLibraryEntry, 0, len(entries))
	deleted := false
	for _, entry := range entries {
		if entry.ID == id {
			deleted = true
			continue
		}
		nextEntries = append(nextEntries, entry)
	}
	if !deleted {
		return fmt.Errorf("custom prompt not found")
	}
	if err := a.saveCustomPromptEntries(nextEntries); err != nil {
		return err
	}

	state, err := a.loadPromptAssistantState()
	if err != nil {
		return err
	}
	state.FavoriteIDs = removeStringValue(state.FavoriteIDs, id)
	state.RecentIDs = removeStringValue(state.RecentIDs, id)
	return a.savePromptAssistantState(state)
}

func (a *App) loadPromptAssistantState() (PromptAssistantState, error) {
	var state PromptAssistantState
	data, err := os.ReadFile(a.promptAssistantStateFile())
	if err != nil {
		if os.IsNotExist(err) {
			return PromptAssistantState{}, nil
		}
		return PromptAssistantState{}, err
	}
	if len(data) == 0 {
		return PromptAssistantState{}, nil
	}
	if err := json.Unmarshal(data, &state); err != nil {
		return PromptAssistantState{}, err
	}
	return normalizePromptAssistantState(state), nil
}

func (a *App) savePromptAssistantState(state PromptAssistantState) error {
	state = normalizePromptAssistantState(state)
	data, _ := json.MarshalIndent(state, "", "  ")
	return os.WriteFile(a.promptAssistantStateFile(), data, 0644)
}

func (a *App) GetPromptAssistantState() (PromptAssistantState, error) {
	return a.loadPromptAssistantState()
}

func (a *App) SavePromptAssistantState(state PromptAssistantState) (PromptAssistantState, error) {
	state = normalizePromptAssistantState(state)
	if err := a.savePromptAssistantState(state); err != nil {
		return PromptAssistantState{}, err
	}
	return state, nil
}

func (a *App) loadPromptToolLinks() ([]PromptToolLink, error) {
	var links []PromptToolLink
	data, err := os.ReadFile(a.promptToolLinksFile())
	if err != nil {
		return []PromptToolLink{}, nil
	}
	json.Unmarshal(data, &links)
	return links, nil
}

func (a *App) savePromptToolLinks(links []PromptToolLink) error {
	data, _ := json.MarshalIndent(links, "", "  ")
	return os.WriteFile(a.promptToolLinksFile(), data, 0644)
}

func (a *App) GetPromptToolLinks() ([]PromptToolLink, error) {
	return a.loadPromptToolLinks()
}

func (a *App) AddPromptToolLink(link PromptToolLink) (PromptToolLink, error) {
	links, _ := a.loadPromptToolLinks()
	link.ID = uuid.New().String()
	links = append(links, link)
	err := a.savePromptToolLinks(links)
	return link, err
}

func (a *App) UpdatePromptToolLink(id string, link PromptToolLink) error {
	links, _ := a.loadPromptToolLinks()
	updated := false
	for i, l := range links {
		if l.ID == id {
			links[i].Name = link.Name
			links[i].URL = link.URL
			links[i].Icon = link.Icon
			updated = true
			break
		}
	}
	if !updated {
		return fmt.Errorf("link not found")
	}
	return a.savePromptToolLinks(links)
}

func (a *App) DeletePromptToolLink(id string) error {
	links, _ := a.loadPromptToolLinks()
	newLinks := []PromptToolLink{}
	for _, l := range links {
		if l.ID != id {
			newLinks = append(newLinks, l)
		}
	}
	return a.savePromptToolLinks(newLinks)
}

func (a *App) loadPromptTemplates() ([]PromptTemplate, error) {
	var templates []PromptTemplate
	data, err := os.ReadFile(a.promptTemplatesFile())
	if err != nil {
		return []PromptTemplate{}, nil
	}
	json.Unmarshal(data, &templates)
	return templates, nil
}

func (a *App) savePromptTemplates(templates []PromptTemplate) error {
	data, _ := json.MarshalIndent(templates, "", "  ")
	return os.WriteFile(a.promptTemplatesFile(), data, 0644)
}

func (a *App) GetPromptTemplates() ([]PromptTemplate, error) {
	return a.loadPromptTemplates()
}

func (a *App) AddPromptTemplate(template PromptTemplate) (PromptTemplate, error) {
	templates, _ := a.loadPromptTemplates()
	template.ID = uuid.New().String()
	template.CreatedAt = time.Now().Format(time.RFC3339)
	templates = append(templates, template)
	err := a.savePromptTemplates(templates)
	return template, err
}

func (a *App) UpdatePromptTemplate(id string, template PromptTemplate) error {
	templates, _ := a.loadPromptTemplates()
	updated := false
	for i, t := range templates {
		if t.ID == id {
			templates[i].Name = template.Name
			templates[i].Content = template.Content
			templates[i].Type = template.Type
			templates[i].Category = template.Category
			templates[i].SourcePath = template.SourcePath
			updated = true
			break
		}
	}
	if !updated {
		return fmt.Errorf("template not found")
	}
	return a.savePromptTemplates(templates)
}

func (a *App) DeletePromptTemplate(id string) error {
	templates, _ := a.loadPromptTemplates()
	newTemplates := []PromptTemplate{}
	for _, t := range templates {
		if t.ID != id {
			newTemplates = append(newTemplates, t)
		}
	}
	return a.savePromptTemplates(newTemplates)
}

// ImportCustomPromptsFromText opens a file dialog, reads a CSV/TXT prompt file, and imports entries into custom prompts.
func (a *App) ImportCustomPromptsFromText() (int, error) {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择提示词库文件",
		Filters: []runtime.FileFilter{
			{DisplayName: "Supported Files", Pattern: "*.txt;*.csv;*.json"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
	if err != nil {
		return 0, err
	}
	if filePath == "" {
		return 0, nil // User cancelled
	}

	f, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	var newEntries []PromptLibraryEntry
	importedSource := "外部导入"

	if strings.HasSuffix(strings.ToLower(filePath), ".json") {
		var jsonEntries []PromptLibraryEntry
		if err := json.NewDecoder(f).Decode(&jsonEntries); err != nil {
			return 0, fmt.Errorf("failed to parse json: %w", err)
		}
		for _, entry := range jsonEntries {
			entry.Source = importedSource
			if entry.Category == "" {
				entry.Category = "default"
			}
			if entry.Subcategory == "" {
				entry.Subcategory = "default"
			}
			if entry.Scope == "" {
				entry.Scope = "normal"
			}
			newEntries = append(newEntries, entry)
		}
	} else {
		reader := csv.NewReader(f)
		reader.FieldsPerRecord = -1 // Allow variable number of fields
		records, err := reader.ReadAll()
		if err != nil {
			return 0, fmt.Errorf("failed to read csv data: %w", err)
		}
		for _, row := range records {
			if len(row) < 4 {
				continue
			}
			en := strings.TrimSpace(row[0])
			lastCol := strings.TrimSpace(row[3])
			parts := strings.Split(lastCol, ",")
			zh := ""
			if len(parts) > 0 {
				zh = strings.TrimSpace(parts[len(parts)-1])
			}
			newEntries = append(newEntries, PromptLibraryEntry{
				Source:      importedSource,
				Category:    "default",
				Subcategory: "default",
				Scope:       "normal",
				TextEN:      en,
				TextZH:      zh,
			})
		}
	}


	customEntries, err := a.loadCustomPromptEntries()
	if err != nil {
		customEntries = []PromptLibraryEntry{}
	}
	custSet := make(map[string]bool)
	for _, item := range customEntries {
		if en := normalizePromptTextKey(item.TextEN); en != "" {
			custSet[en] = true
		}
		if zh := normalizePromptTextKey(item.TextZH); zh != "" {
			custSet[zh] = true
		}
	}

	addedCount := 0
	var finalEntries []PromptLibraryEntry

	for _, entry := range newEntries {
		if entry.TextEN == "" && entry.TextZH == "" {
			continue
		}

		normEN := normalizePromptTextKey(entry.TextEN)
		normZH := normalizePromptTextKey(entry.TextZH)


		if (normEN != "" && custSet[normEN]) || (normZH != "" && custSet[normZH]) {
			continue
		}

		if normEN != "" {
			custSet[normEN] = true
		}
		if normZH != "" {
			custSet[normZH] = true
		}

		entry.ID = uuid.New().String()
		entry = normalizePromptLibraryEntry(entry)
		finalEntries = append(finalEntries, entry)
		addedCount++
	}

	if addedCount > 0 {
		customEntries = append(finalEntries, customEntries...)
		if err := a.saveCustomPromptEntries(customEntries); err != nil {
			return 0, fmt.Errorf("failed to save custom prompts: %w", err)
		}
	}

	return addedCount, nil
}

// ClearExternalImports removes all prompt entries with source '外部导入' from custom prompts.
func (a *App) ClearExternalImports() error {
	customEntries, err := a.loadCustomPromptEntries()
	if err != nil {
		return err
	}

	var filtered []PromptLibraryEntry
	for _, entry := range customEntries {
		if entry.Source != "外部导入" {
			filtered = append(filtered, entry)
		}
	}

	return a.saveCustomPromptEntries(filtered)
}
