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
	Tokens      []string
	RegexTokens []RegexToken
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
			if regularToken := l.testRegularToken(token); regularToken != "" {
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

func (l *Lexer) tokensFromFile(filePath string) ([][]string, error) {
	var tokens [][]string
	file, err := os.Open(filePath)
	if err != nil {
		return tokens, err
	}
	reader := bufio.NewReader(file)
	regexString := l.fullRegex()
	for {
		line, err := reader.ReadString('\n')
		x := regexp.MustCompile(regexString)
		array := x.FindAllStringSubmatch(line, -1)

		var tokensOnLine []string
		for _, i := range array {
			for _, j := range i {
				tokensOnLine = append(tokensOnLine, j)
			}
		}
		tokens = append(tokens, tokensOnLine)
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}
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
