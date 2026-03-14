package services

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"
)

const (
	GreetingResponse       = "Hello, ask me about Chintak's projects, skills, or experience."
	AcknowledgmentResponse = "Thanks. Ask me about another project, skill, or part of Chintak's experience."
	PortfolioOnlyResponse  = "I can only answer questions about Chintak Joshi's portfolio. Ask about his projects, skills, experience, education, achievements, or links listed there."
)

var (
	wordPattern       = regexp.MustCompile(`[a-z0-9]+`)
	arithmeticPattern = regexp.MustCompile(`\d+\s*[\+\-\*/]\s*\d+`)
	canonicalPattern  = regexp.MustCompile(`[^a-z0-9]+`)
)

type GuardrailDecision struct {
	DirectResponse   string
	SystemPrompt     string
	UserPrompt       string
	KnowledgeContext string
}

const sessionContextTTL = 30 * time.Minute

type sessionContext struct {
	knowledgeContext string
	updatedAt        time.Time
}

var (
	sessionContextMu    sync.RWMutex
	sessionContextStore = make(map[string]sessionContext)
)

func ApplyGuardrails(userMessage string) GuardrailDecision {
	return ApplyGuardrailsWithSession("", userMessage)
}

func ApplyGuardrailsWithSession(sessionID, userMessage string) GuardrailDecision {
	message := normalizeForMatching(userMessage)
	if message == "" {
		return GuardrailDecision{DirectResponse: PortfolioOnlyResponse}
	}

	if isGreetingMessage(message) {
		return GuardrailDecision{DirectResponse: GreetingResponse}
	}

	if isAcknowledgmentMessage(message) {
		return GuardrailDecision{DirectResponse: AcknowledgmentResponse}
	}

	if isPromptInjectionAttempt(message) {
		return GuardrailDecision{DirectResponse: PortfolioOnlyResponse}
	}

	intents := detectIntents(message)
	if isFollowUpMessage(message) {
		context, ok := getSessionContext(sessionID)
		if ok && !isUnsupportedTask(message, intents) {
			return rememberDecision(sessionID, GuardrailDecision{
				SystemPrompt:     GetScopedContextPrompt(context),
				UserPrompt:       EnhanceFollowUpUserMessage(userMessage),
				KnowledgeContext: context,
			})
		}
	}

	context, allow := buildScopedKnowledgeContext(message, intents)
	if !allow || isUnsupportedTask(message, intents) {
		return GuardrailDecision{DirectResponse: PortfolioOnlyResponse}
	}

	return rememberDecision(sessionID, GuardrailDecision{
		SystemPrompt:     GetScopedContextPrompt(context),
		UserPrompt:       EnhanceUserMessage(userMessage),
		KnowledgeContext: context,
	})
}

func buildScopedKnowledgeContext(message string, intents map[string]bool) (string, bool) {
	personalInfo := GetPersonalInfo()
	projects := GetProjects()
	portfolioReference := hasPortfolioReference(message)

	var sections []string

	if shouldIncludeProfile(message, intents, portfolioReference) {
		sections = append(sections, buildProfileSection(personalInfo))
	}

	if intents["education"] {
		sections = append(sections, buildEducationSection(personalInfo))
	}

	if intents["achievements"] {
		sections = append(sections, buildAchievementsSection(personalInfo))
	}

	matchingSkillCategories := findMatchingSkillCategories(message, personalInfo.Skills)
	if intents["skills"] || len(matchingSkillCategories) > 0 {
		sections = append(sections, buildScopedSkillsSection(personalInfo.Skills, matchingSkillCategories))
	}

	matchingProjects := findMatchingProjects(message, projects)
	if intents["projects"] || len(matchingProjects) > 0 {
		sections = append(sections, buildScopedProjectsSection(projects, matchingProjects))
	}

	if intents["contact"] || referencesKnownLink(message, personalInfo.Links) {
		sections = append(sections, buildLinksSection(personalInfo.Links))
	}

	score := 0
	if countMatchedIntents(intents) > 0 {
		score += 2
	}
	if portfolioReference {
		score++
	}
	if len(matchingProjects) > 0 {
		score += 2
	}
	if len(matchingSkillCategories) > 0 {
		score++
	}
	if referencesPersonalInfo(message, personalInfo) || referencesKnownLink(message, personalInfo.Links) {
		score++
	}

	contextSections := compactSections(sections)
	if len(contextSections) == 0 || score < 2 {
		return "", false
	}

	return strings.Join(contextSections, "\n\n"), true
}

func shouldIncludeProfile(message string, intents map[string]bool, portfolioReference bool) bool {
	if intents["about"] || portfolioReference {
		return true
	}

	return containsAny(message, []string{
		"current role", "current job", "work", "working", "experience",
		"title", "name", "who is", "who are", "united rentals",
	})
}

func buildProfileSection(info PersonalInfo) string {
	return fmt.Sprintf(`PROFILE:
- Name: %s
- Title: %s
- Experience: %s
- Current Role: %s`, info.Name, info.Title, info.Experience, info.CurrentRole)
}

func buildEducationSection(info PersonalInfo) string {
	return "EDUCATION:\n- " + strings.Join(info.Education, "\n- ")
}

func buildAchievementsSection(info PersonalInfo) string {
	return "ACHIEVEMENTS:\n- " + strings.Join(info.Achievements, "\n- ")
}

func buildScopedSkillsSection(skills map[string][]string, categories []string) string {
	if len(categories) == 0 {
		return buildSkillsSection(skills)
	}

	var b strings.Builder
	b.WriteString("SKILLS:\n")
	for _, category := range categories {
		b.WriteString(fmt.Sprintf("- %s: %s\n", category, strings.Join(skills[category], ", ")))
	}

	return strings.TrimSpace(b.String())
}

func buildScopedProjectsSection(allProjects []Project, matchingProjects []Project) string {
	projectsToUse := matchingProjects
	if len(projectsToUse) == 0 {
		projectsToUse = allProjects
	}

	return buildProjectsSection(projectsToUse)
}

func findMatchingSkillCategories(message string, skills map[string][]string) []string {
	var matches []string

	for _, category := range sortedStringKeys(skills) {
		categoryLower := strings.ToLower(category)
		if strings.Contains(message, categoryLower) {
			matches = append(matches, category)
			continue
		}

		for _, skill := range skills[category] {
			if strings.Contains(message, strings.ToLower(skill)) {
				matches = append(matches, category)
				break
			}
		}
	}

	return compactStrings(matches)
}

func findMatchingProjects(message string, projects []Project) []Project {
	messageTokens := tokenSet(message)
	matches := make([]Project, 0)

	for _, project := range projects {
		projectText := strings.ToLower(strings.Join([]string{
			project.Name,
			project.Description,
			project.Category,
			strings.Join(project.Technologies, " "),
			strings.Join(project.Features, " "),
		}, " "))

		if strings.Contains(message, strings.ToLower(project.Name)) {
			matches = append(matches, project)
			continue
		}

		overlap := overlapScore(messageTokens, tokenSet(projectText))
		if overlap >= 2 {
			matches = append(matches, project)
		}
	}

	return matches
}

func referencesPersonalInfo(message string, info PersonalInfo) bool {
	fields := []string{
		info.Name,
		info.Title,
		info.Experience,
		info.CurrentRole,
	}

	for _, field := range fields {
		if strings.Contains(message, strings.ToLower(field)) {
			return true
		}
	}

	return false
}

func referencesKnownLink(message string, links map[string]string) bool {
	for platform := range links {
		if strings.Contains(message, strings.ToLower(platform)) {
			return true
		}
	}

	return false
}

func hasPortfolioReference(message string) bool {
	return containsAny(message, []string{
		"you", "your", "yourself", "chintak", "joshi", "he", "his",
	})
}

func isGreetingMessage(message string) bool {
	greetings := []string{
		"hi", "hello", "hey", "good morning", "good afternoon", "good evening",
	}

	return slices.Contains(greetings, canonicalMessage(message))
}

func isAcknowledgmentMessage(message string) bool {
	acknowledgments := []string{
		"thanks",
		"thank you",
		"looks good",
		"looks great",
		"sounds good",
		"impressive",
		"nice",
		"cool",
		"awesome",
		"great",
		"ok",
		"great work",
		"well done",
		"nice",
		"good job",
		"awesome",
		"fantastic",
		"amazing",
		"excellent",
		"wonderful",
		"superb",
		"brilliant",
		"outstanding",
		"fine",
		"acceptable",
	}

	return slices.Contains(acknowledgments, canonicalMessage(message))
}

func isFollowUpMessage(message string) bool {
	return containsAny(message, []string{
		"tell me more",
		"tell me in detail",
		"can you tell me in detail",
		"more detail",
		"more details",
		"in detail",
		"in more detail",
		"elaborate",
		"expand on that",
		"explain more",
		"more about that",
		"more about it",
		"what about that",
		"what about it",
	})
}

func isPromptInjectionAttempt(message string) bool {
	return containsAny(message, []string{
		"ignore previous instructions",
		"ignore all previous",
		"reveal system prompt",
		"show system prompt",
		"developer message",
		"act as",
		"jailbreak",
		"override instructions",
	})
}

func isUnsupportedTask(message string, intents map[string]bool) bool {
	taskPhrases := []string{
		"write a", "write an", "create a", "generate a", "build a", "implement",
		"solve", "calculate", "compute", "code", "function", "script", "program",
		"algorithm", "equation", "homework", "translate", "poem", "story", "recipe",
	}

	if arithmeticPattern.MatchString(message) {
		return true
	}

	if containsAny(message, taskPhrases) && countMatchedIntents(intents) == 0 {
		return true
	}

	return false
}

func normalizeForMatching(message string) string {
	return strings.TrimSpace(strings.ToLower(message))
}

func canonicalMessage(message string) string {
	normalized := normalizeForMatching(message)
	normalized = canonicalPattern.ReplaceAllString(normalized, " ")
	return strings.Join(strings.Fields(normalized), " ")
}

func compactSections(sections []string) []string {
	return compactStrings(sections)
}

func compactStrings(values []string) []string {
	var compacted []string
	seen := make(map[string]bool)

	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}

		seen[value] = true
		compacted = append(compacted, value)
	}

	return compacted
}

func tokenSet(text string) map[string]bool {
	tokens := make(map[string]bool)
	for _, token := range wordPattern.FindAllString(strings.ToLower(text), -1) {
		if len(token) < 3 || isStopWord(token) {
			continue
		}
		tokens[token] = true
	}

	return tokens
}

func overlapScore(a, b map[string]bool) int {
	score := 0
	for token := range a {
		if b[token] {
			score++
		}
	}
	return score
}

func isStopWord(token string) bool {
	stopWords := map[string]bool{
		"the": true, "and": true, "for": true, "are": true, "with": true,
		"what": true, "when": true, "where": true, "have": true, "from": true,
		"that": true, "this": true, "about": true, "your": true, "you": true,
		"his": true, "her": true, "who": true, "how": true, "can": true,
	}

	return stopWords[token]
}

func EnhanceFollowUpUserMessage(userMessage string) string {
	return userMessage + " [STRICT: This is a follow-up to the previously discussed portfolio topic. Use only the provided portfolio information, stay on the same topic, and give more detail only from that information.]"
}

func rememberDecision(sessionID string, decision GuardrailDecision) GuardrailDecision {
	if sessionID == "" || decision.KnowledgeContext == "" {
		return decision
	}

	sessionContextMu.Lock()
	sessionContextStore[sessionID] = sessionContext{
		knowledgeContext: decision.KnowledgeContext,
		updatedAt:        time.Now(),
	}
	sessionContextMu.Unlock()

	return decision
}

func getSessionContext(sessionID string) (string, bool) {
	if sessionID == "" {
		return "", false
	}

	sessionContextMu.RLock()
	context, ok := sessionContextStore[sessionID]
	sessionContextMu.RUnlock()
	if !ok {
		return "", false
	}

	if time.Since(context.updatedAt) > sessionContextTTL {
		sessionContextMu.Lock()
		delete(sessionContextStore, sessionID)
		sessionContextMu.Unlock()
		return "", false
	}

	return context.knowledgeContext, true
}
