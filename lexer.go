package main

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode"
)

type Lexer struct {
	Tokens       []string
	RegexTokens  []RegexToken
	TokenTable   TokenList
	CurrentToken int
}

type TokenState struct {
	Tokens   TokenList
	curToken int
}

type TokenList []Token
type Token struct {
	Name  string
	Value string
	Line  uint
}
type RegexTokenList []RegexToken
type RegexToken struct {
	Name  string
	Regex string
}

func newLexerFromFile(fileName string) (Lexer, error) {
	var lexer Lexer
	file, err := os.Open(fileName)
	if err != nil {
		return lexer, err
	}
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return lexer, err
			}
		}
		lineElements := strings.Split(string(line), " ")
		lineElements = lineElements[1:]
		switch line[0] {
		case '#':
			lexer.Tokens = append(lexer.Tokens, lineElements...)
		case '%':
			for i := 0; i < len(lineElements)-1; i += 2 {
				var regexToken RegexToken
				regexToken.Name = lineElements[i]
				regexToken.Regex = lineElements[i+1]
				lexer.RegexTokens = append(lexer.RegexTokens, regexToken)
			}
		}
	}
	return lexer, nil
}

func (l *Lexer) testTokensOfFile(filePath string) (TokenList, error) {
	var tokenList TokenList
	fileTokens, err := l.tokensFromFile(filePath)
	if err != nil {
		return tokenList, err
	}
	for i, line := range fileTokens {
		for _, token := range line {
			var t Token
			if token == "INDENT" {
				t.Line = uint(i) + 1
				t.Name = "INDENT"
				t.Value = "INDENT"
				tokenList = append(tokenList, t)
				t = Token{}
			} else if token == "DEDENT" {
				t.Line = uint(i) + 1
				t.Name = "DEDENT"
				t.Value = "DEDENT"
				tokenList = append(tokenList, t)
				t = Token{}
			} else if token == "EOF" {
				t.Line = uint(i) + 1
				t.Name = "EOF"
				t.Value = "EOF"
				tokenList = append(tokenList, t)
				t = Token{}
			} else if regularToken := l.testRegularToken(token); regularToken != "" {
				t.Line = uint(i) + 1
				t.Name = regularToken
				t.Value = regularToken
				tokenList = append(tokenList, t)
				t = Token{}
			} else if regexToken := l.testRegexToken(token); regexToken != "" {
				t.Line = uint(i) + 1
				t.Name = regexToken
				t.Value = token
				tokenList = append(tokenList, t)
				t = Token{}
			} else {
				t.Line = uint(i) + 1
				t.Name = "INVALID"
				t.Value = token
				tokenList = append(tokenList, t)
				return tokenList, nil
			}
		}
	}
	return tokenList, nil
}
func AppendNewLineToFile(path string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("\n")
	if err != nil {
		return err
	}
	return nil
}
func (l *Lexer) tokensFromFile(filePath string) ([][]string, error) {
	var tokens [][]string
	file, err := os.Open(filePath)
	if err != nil {
		return tokens, err
	}
	reader := bufio.NewReader(file)
	regexString := l.fullRegex()
	identSpaces := 0
	previousLineIdentLevel := 0
	currentLineIdentLevel := 0
	for {
		var tokensOnLine []string
		line, err := reader.ReadString('\n')
		currentLineIdentLevel = defineLineIdentLevel(line)
		if currentLineIdentLevel > previousLineIdentLevel && identSpaces == 0 {
			tokensOnLine = append(tokensOnLine, "INDENT")
			identSpaces = currentLineIdentLevel
			previousLineIdentLevel = currentLineIdentLevel
		} else {
			if currentLineIdentLevel > previousLineIdentLevel {
				indentCount := (currentLineIdentLevel - previousLineIdentLevel) / identSpaces
				for i := 0; i < indentCount; i++ {
					tokensOnLine = append(tokensOnLine, "INDENT")
				}
				previousLineIdentLevel = currentLineIdentLevel
			} else if currentLineIdentLevel < previousLineIdentLevel {
				dedentCount := (previousLineIdentLevel - currentLineIdentLevel) / identSpaces
				for i := 0; i < dedentCount; i++ {
					tokensOnLine = append(tokensOnLine, "DEDENT")
				}
				previousLineIdentLevel = currentLineIdentLevel
			}
		}
		x := regexp.MustCompile(regexString)
		array := x.FindAllStringSubmatch(line, -1)
		for _, i := range array {
			for _, j := range i {
				tokensOnLine = append(tokensOnLine, j)
			}
		}
		if err != nil {
			if err == io.EOF {
				lastElement := tokensOnLine[len(tokensOnLine)-1]
				if lastElement != "DEDENT" && currentLineIdentLevel != 0 && identSpaces != 0 {
					dedentCount := currentLineIdentLevel / identSpaces
					for i := 0; i < dedentCount; i++ {
						tokensOnLine = append(tokensOnLine, "DEDENT")
					}
				}
				tokensOnLine = append(tokensOnLine, "EOF")
				tokens = append(tokens, tokensOnLine)
				break
			}
			break
		}
		tokens = append(tokens, tokensOnLine)
	}
	return tokens, nil
}

func (l *Lexer) testRegularToken(t string) string {
	for _, token := range l.Tokens {
		//HACK FIX
		if token[0] == '\\' {
			re := regexp.MustCompile(`\\`)
			newToken := re.ReplaceAllString(token, "")
			newT := re.ReplaceAllString(t, "")
			if newToken == newT {
				return newT
			}
		}
		if token == t {
			return token
		}
	}
	return ""
}

func (l *Lexer) testRegexToken(element string) string {
	for _, token := range l.RegexTokens {
		newRegex := "^"
		newRegex += token.Regex + "$"
		if match, err := regexp.MatchString(newRegex, element); err == nil && match == true {
			return token.Name
		} else if err != nil {
			return ""
		}
	}
	return ""
}

func (l *Lexer) next() Token {
	l.CurrentToken++
	return l.TokenTable[l.CurrentToken-1]
}

func (l *Lexer) back() Token {
	l.CurrentToken--
	return l.TokenTable[l.CurrentToken-1]
}

func (l *Lexer) fullRegex() string {
	regexString := ""
	for _, token := range l.Tokens {
		regexString += token + `|`
	}
	for _, regexToken := range l.RegexTokens {
		regexString += regexToken.Regex + `|`
	}
	return regexString + `[^\s.*]+`
}

func removeSpace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func defineLineIdentLevel(s string) int {
	counter := 0
	for counter < len(s)-1 {
		if s[counter] == ' ' {
			counter++
		} else {
			break
		}
	}
	return counter
}

func (l *Lexer) isEOF() bool {
	return l.TokenTable[l.CurrentToken].Name == "EOF"
}
