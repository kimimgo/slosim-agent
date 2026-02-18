package prompt

import (
	"os"
	"strings"
	"testing"

	"github.com/opencode-ai/opencode/internal/llm/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSloshingCoderPrompt_Full(t *testing.T) {
	os.Unsetenv("SLOSIM_ABLATION")
	prompt := SloshingCoderPrompt(models.ModelProvider(""))

	assert.Contains(t, prompt, "슬로싱(Sloshing) 해석 전문")
	assert.Contains(t, prompt, "f₁ = (1/2π)")
	assert.Contains(t, prompt, "LNG 탱크")
	assert.Contains(t, prompt, "xml_generator → XML")
	assert.Contains(t, prompt, "simulations/{case_name}/")
}

func TestSloshingCoderPrompt_NoDomain(t *testing.T) {
	os.Setenv("SLOSIM_ABLATION", "no-domain")
	defer os.Unsetenv("SLOSIM_ABLATION")

	prompt := SloshingCoderPrompt(models.ModelProvider(""))

	// Should keep: role, rules, tool order, path rules
	assert.Contains(t, prompt, "슬로싱(Sloshing) 해석 전문")
	assert.Contains(t, prompt, "xml_generator → XML")
	assert.Contains(t, prompt, "simulations/{case_name}/")

	// Should remove: domain knowledge, tank presets, parameter rules
	assert.NotContains(t, prompt, "f₁ = (1/2π)")
	assert.NotContains(t, prompt, "LNG 탱크")
	assert.NotContains(t, prompt, "dp = min(L,W,H)/50")
}

func TestSloshingCoderPrompt_NoRules(t *testing.T) {
	os.Setenv("SLOSIM_ABLATION", "no-rules")
	defer os.Unsetenv("SLOSIM_ABLATION")

	prompt := SloshingCoderPrompt(models.ModelProvider(""))

	// Should keep: role, domain knowledge, features
	assert.Contains(t, prompt, "슬로싱(Sloshing) 해석 전문")
	assert.Contains(t, prompt, "f₁ = (1/2π)")
	assert.Contains(t, prompt, "LNG 탱크")

	// Should remove: absolute rules, tool call order, path rules
	assert.NotContains(t, prompt, "절대 규칙")
	assert.NotContains(t, prompt, "xml_generator → XML")
	assert.NotContains(t, prompt, "simulations/{case_name}/")
}

func TestSloshingCoderPrompt_Generic(t *testing.T) {
	os.Setenv("SLOSIM_ABLATION", "generic")
	defer os.Unsetenv("SLOSIM_ABLATION")

	prompt := SloshingCoderPrompt(models.ModelProvider(""))

	assert.Contains(t, prompt, "helpful AI assistant")
	assert.Contains(t, prompt, "DualSPHysics")

	// Should be very short
	assert.Less(t, len(prompt), 200)

	// Should NOT contain any domain knowledge
	assert.NotContains(t, prompt, "f₁ = (1/2π)")
	assert.NotContains(t, prompt, "절대 규칙")
}

func TestSloshingCoderPrompt_DefaultIsFull(t *testing.T) {
	os.Unsetenv("SLOSIM_ABLATION")
	full := SloshingCoderPrompt(models.ModelProvider(""))

	os.Setenv("SLOSIM_ABLATION", "full")
	defer os.Unsetenv("SLOSIM_ABLATION")
	explicit := SloshingCoderPrompt(models.ModelProvider(""))

	assert.Equal(t, full, explicit)
}

func TestSloshingCoderPrompt_UnknownAblationDefaultsToFull(t *testing.T) {
	os.Setenv("SLOSIM_ABLATION", "unknown-mode")
	defer os.Unsetenv("SLOSIM_ABLATION")

	prompt := SloshingCoderPrompt(models.ModelProvider(""))
	assert.Contains(t, prompt, "f₁ = (1/2π)")
}

func TestGetAblationMode(t *testing.T) {
	os.Unsetenv("SLOSIM_ABLATION")
	assert.Equal(t, "full", GetAblationMode())

	os.Setenv("SLOSIM_ABLATION", "no-domain")
	defer os.Unsetenv("SLOSIM_ABLATION")
	assert.Equal(t, "no-domain", GetAblationMode())
}

func TestPromptSections_NoDuplication(t *testing.T) {
	// Verify the full prompt equals concatenation of all sections
	expected := promptRole +
		promptAbsoluteRules +
		promptToolCallOrder +
		promptParameterRules +
		promptTankPresets +
		promptDomainKnowledge +
		promptToolDetailRules +
		promptFolderRules +
		promptFeatures +
		promptParaView +
		promptConstraints

	assert.Equal(t, expected, sloshingSystemPrompt)
}

func TestPromptSections_TokenCounts(t *testing.T) {
	// Rough token count check (chars/4 approximation)
	// Full prompt should be reasonable for Qwen3 32B context
	fullLen := len(sloshingSystemPrompt)
	require.Greater(t, fullLen, 1000, "full prompt too short")

	// Generic should be < 200 chars
	require.Less(t, len(sloshingGenericPrompt), 200)

	// NO-DOMAIN should be shorter than FULL (removed 3 sections)
	noDomainLen := len(sloshingNoDomainPrompt)
	assert.Less(t, noDomainLen, fullLen)

	// NO-RULES should be shorter than FULL (removed 4 sections)
	noRulesLen := len(sloshingNoRulesPrompt)
	assert.Less(t, noRulesLen, fullLen)

	// Log sizes for paper reporting
	t.Logf("Prompt sizes (chars): FULL=%d, NO-DOMAIN=%d, NO-RULES=%d, GENERIC=%d",
		fullLen, noDomainLen, noRulesLen, len(sloshingGenericPrompt))
	t.Logf("Prompt sizes (approx tokens): FULL≈%d, NO-DOMAIN≈%d, NO-RULES≈%d, GENERIC≈%d",
		fullLen/4, noDomainLen/4, noRulesLen/4, len(sloshingGenericPrompt)/4)
}

func TestPromptContentIntegrity(t *testing.T) {
	// The full prompt should contain all section markers
	markers := []string{
		"# 역할",
		"# 절대 규칙",
		"# Tool 호출 순서",
		"# 파라미터 자동 결정 규칙",
		"# 표준 탱크 치수",
		"# 도메인 지식",
		"# Tool 호출 세부 규칙",
		"# 시뮬레이션 결과 폴더 규칙",
		"# 지원 기능",
		"# ParaView 후처리 도구",
		"# 제약 사항",
	}

	for _, marker := range markers {
		count := strings.Count(sloshingSystemPrompt, marker)
		assert.Equal(t, 1, count, "marker %q should appear exactly once, found %d", marker, count)
	}
}
