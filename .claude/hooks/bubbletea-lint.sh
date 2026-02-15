#!/bin/bash
# PostToolUse hook: BubbleTea TUI 디자인 규칙 검증
# internal/tui/components/ 편집 시 하드코딩 색상, 고정 크기 차단

INPUT=$(cat)
FILE_PATH=$(echo "$INPUT" | jq -r '.tool_input.file_path // .tool_input.filePath // empty')

# TUI 컴포넌트만 대상
if [[ "$FILE_PATH" != */internal/tui/components/* ]]; then exit 0; fi
if [[ "$FILE_PATH" != *.go ]]; then exit 0; fi

ERRORS=()

# 1. 하드코딩 색상 검출 (lipgloss.Color("숫자") 패턴)
#    theme_test.go, *_test.go는 제외
if [[ "$FILE_PATH" != *_test.go ]]; then
  if grep -nE 'lipgloss\.Color\("[0-9]' "$FILE_PATH" | grep -v '//.*lint:ignore' > /dev/null 2>&1; then
    LINES=$(grep -nE 'lipgloss\.Color\("[0-9]' "$FILE_PATH" | head -3)
    ERRORS+=("하드코딩 색상 금지 — theme.CurrentTheme() 사용. 위반:\n$LINES")
  fi
fi

# 2. 고정 Width/Height 검출 (3자리 이상 숫자)
if grep -nE '\.(Width|Height)\([0-9]{3,}\)' "$FILE_PATH" > /dev/null 2>&1; then
  LINES=$(grep -nE '\.(Width|Height)\([0-9]{3,}\)' "$FILE_PATH" | head -3)
  ERRORS+=("고정 크기 금지 — lipgloss.Width()/msg.Width 동적 계산 사용. 위반:\n$LINES")
fi

# 에러 보고
if [ ${#ERRORS[@]} -gt 0 ]; then
  echo "BubbleTea 디자인 규칙 위반:" >&2
  for err in "${ERRORS[@]}"; do
    echo -e "  - $err" >&2
  done
  echo "" >&2
  echo "참조: .claude/skills/go-tui-design/SKILL.md" >&2
  exit 2  # 차단
fi

exit 0
