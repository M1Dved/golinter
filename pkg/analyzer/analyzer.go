package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "golinter",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var (
	disableRules  string
	extraKeywords string
	configPath    string
)

func init() {
	Analyzer.Flags.StringVar(&configPath, "config", "",
		"path to config file (default: .golinter.yml)")
	Analyzer.Flags.StringVar(&disableRules, "disable", "",
		"comma-separated rules to disable: lowercase,english,special,sensitive")
	Analyzer.Flags.StringVar(&extraKeywords, "extra-keywords", "",
		"comma-separated additional sensitive keywords")
}

var logMethods = map[string]bool{
	"Debug": true, "Info": true, "Warn": true, "Error": true,
	"DPanic": true, "Fatal": true, "Panic": true,
	"Debugf": true, "Infof": true, "Warnf": true, "Errorf": true,
	"DPanicf": true, "Fatalf": true, "Panicf": true,
	"Debugw": true, "Infow": true, "Warnw": true, "Errorw": true,
}

var logPackages = map[string]bool{
	"log": true, "slog": true, "zap": true,
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	disabled, extra := buildSettings()

	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return
		}

		if _, ok := isLogCall(call); !ok {
			return
		}

		msg, pos, ok := extractMessage(call)
		if !ok {
			return
		}

		checkRules(pass, msg, pos, disabled, extra)
	})

	return nil, nil
}

func buildSettings() (map[string]bool, []string) {
	disabled := make(map[string]bool)
	var extra []string

	cfg := tryLoadConfig()
	if cfg != nil {
		disabled = cfg.DisabledRules()
		for _, kw := range cfg.ExtraKeywords {
			extra = append(extra, strings.ToLower(strings.TrimSpace(kw)))
		}
	}

	if disableRules != "" {
		for _, rule := range strings.Split(disableRules, ",") {
			disabled[strings.TrimSpace(rule)] = true
		}
	}

	if extraKeywords != "" {
		for _, kw := range strings.Split(extraKeywords, ",") {
			kw = strings.TrimSpace(kw)
			if kw != "" {
				extra = append(extra, strings.ToLower(kw))
			}
		}
	}

	return disabled, extra
}

func tryLoadConfig() *Config {
	path := configPath
	if path == "" {
		path = ".golinter.yml"
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		return nil
	}
	return cfg
}

func isLogCall(call *ast.CallExpr) (string, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}

	methodName := sel.Sel.Name
	if !logMethods[methodName] {
		return "", false
	}

	if ident, ok := sel.X.(*ast.Ident); ok {
		if logPackages[ident.Name] {
			return methodName, true
		}
	}

	return methodName, true
}

func extractMessage(call *ast.CallExpr) (string, token.Pos, bool) {
	if len(call.Args) == 0 {
		return "", 0, false
	}

	firstArg := call.Args[0]

	if lit, ok := firstArg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		return trimQuotes(lit.Value), lit.Pos(), true
	}

	if binExpr, ok := firstArg.(*ast.BinaryExpr); ok && binExpr.Op == token.ADD {
		if lit, ok := binExpr.X.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			return trimQuotes(lit.Value), lit.Pos(), true
		}
	}

	return "", 0, false
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '`' && s[len(s)-1] == '`') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func checkRules(pass *analysis.Pass, msg string, pos token.Pos, disabled map[string]bool, extraKw []string) {
	if !disabled["lowercase"] && !checkLowercase(msg) {
		pass.Reportf(pos, "log message should start with a lowercase letter")
	}

	if !disabled["english"] && !checkEnglishOnly(msg) {
		pass.Reportf(pos, "log message should be in English only")
	}

	if !disabled["special"] && !checkNoSpecialChars(msg) {
		pass.Reportf(pos, "log message should not contain special characters or emoji")
	}

	if !disabled["sensitive"] && !checkNoSensitiveData(msg, extraKw) {
		pass.Reportf(pos, "log message may contain sensitive data")
	}
}
