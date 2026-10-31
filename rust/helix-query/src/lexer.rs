//! Helix Query Engine - Lexer Module
//! 
//! This module provides the lexer for the Helix query language.
//! Tokenizes PromQL-like query expressions.

use std::fmt;

/// Token types for the query language
#[derive(Debug, Clone, PartialEq)]
pub enum Token {
    // Literals
    Number(f64),
    String(String),
    Identifier(String),
    
    // Operators
    Plus,
    Minus,
    Star,
    Slash,
    Percent,
    Caret,
    Equal,
    NotEqual,
    LessThan,
    GreaterThan,
    LessThanOrEqual,
    GreaterThanOrEqual,
    RegexMatch,
    RegexNotMatch,
    
    // Delimiters
    LeftParen,
    RightParen,
    LeftBrace,
    RightBrace,
    LeftBracket,
    RightBracket,
    Comma,
    Dot,
    
    // Keywords
    By,
    Without,
    On,
    Ignoring,
    GroupLeft,
    GroupRight,
    Bool,
    Offset,
    At,
    
    // Special
    Eof,
    Error(String),
}

impl fmt::Display for Token {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Token::Number(n) => write!(f, "Number({})", n),
            Token::String(s) => write!(f, "String({})", s),
            Token::Identifier(s) => write!(f, "Identifier({})", s),
            Token::Plus => write!(f, "+"),
            Token::Minus => write!(f, "-"),
            Token::Star => write!(f, "*"),
            Token::Slash => write!(f, "/"),
            Token::Percent => write!(f, "%"),
            Token::Caret => write!(f, "^"),
            Token::Equal => write!(f, "="),
            Token::NotEqual => write!(f, "!="),
            Token::LessThan => write!(f, "<"),
            Token::GreaterThan => write!(f, ">"),
            Token::LessThanOrEqual => write!(f, "<="),
            Token::GreaterThanOrEqual => write!(f, ">="),
            Token::RegexMatch => write!(f, "=~"),
            Token::RegexNotMatch => write!(f, "!~"),
            Token::LeftParen => write!(f, "("),
            Token::RightParen => write!(f, ")"),
            Token::LeftBrace => write!(f, "{{"),
            Token::RightBrace => write!(f, "}}"),
            Token::LeftBracket => write!(f, "["),
            Token::RightBracket => write!(f, "]"),
            Token::Comma => write!(f, ","),
            Token::Dot => write!(f, "."),
            Token::By => write!(f, "by"),
            Token::Without => write!(f, "without"),
            Token::On => write!(f, "on"),
            Token::Ignoring => write!(f, "ignoring"),
            Token::GroupLeft => write!(f, "group_left"),
            Token::GroupRight => write!(f, "group_right"),
            Token::Bool => write!(f, "bool"),
            Token::Offset => write!(f, "offset"),
            Token::At => write!(f, "@"),
            Token::Eof => write!(f, "EOF"),
            Token::Error(e) => write!(f, "Error({})", e),
        }
    }
}

/// Lexer for the query language
pub struct Lexer {
    input: String,
    pos: usize,
    tokens: Vec<Token>,
}

impl Lexer {
    pub fn new(input: String) -> Self {
        Lexer {
            input,
            pos: 0,
            tokens: Vec::new(),
        }
    }
    
    /// Tokenize the input string
    pub fn tokenize(&mut self) -> Vec<Token> {
        self.tokens.clear();
        
        while self.pos < self.input.len() {
            self.skip_whitespace();
            
            if self.pos >= self.input.len() {
                break;
            }
            
            let c = self.current_char();
            
            match c {
                '+' => { self.pos += 1; self.tokens.push(Token::Plus); }
                '-' => { self.pos += 1; self.tokens.push(Token::Minus); }
                '*' => { self.pos += 1; self.tokens.push(Token::Star); }
                '/' => { self.pos += 1; self.tokens.push(Token::Slash); }
                '%' => { self.pos += 1; self.tokens.push(Token::Percent); }
                '^' => { self.pos += 1; self.tokens.push(Token::Caret); }
                '(' => { self.pos += 1; self.tokens.push(Token::LeftParen); }
                ')' => { self.pos += 1; self.tokens.push(Token::RightParen); }
                '{' => { self.pos += 1; self.tokens.push(Token::LeftBrace); }
                '}' => { self.pos += 1; self.tokens.push(Token::RightBrace); }
                '[' => { self.pos += 1; self.tokens.push(Token::LeftBracket); }
                ']' => { self.pos += 1; self.tokens.push(Token::RightBracket); }
                ',' => { self.pos += 1; self.tokens.push(Token::Comma); }
                '.' => { self.pos += 1; self.tokens.push(Token::Dot); }
                '@' => { self.pos += 1; self.tokens.push(Token::At); }
                '=' => {
                    if self.peek() == '~' {
                        self.pos += 2;
                        self.tokens.push(Token::RegexMatch);
                    } else if self.peek() == '=' {
                        self.pos += 2;
                        self.tokens.push(Token::Equal);
                    } else {
                        self.pos += 1;
                        self.tokens.push(Token::Equal);
                    }
                }
                '!' => {
                    if self.peek() == '=' {
                        self.pos += 2;
                        self.tokens.push(Token::NotEqual);
                    } else if self.peek() == '~' {
                        self.pos += 2;
                        self.tokens.push(Token::RegexNotMatch);
                    } else {
                        self.pos += 1;
                        self.tokens.push(Token::Error("unexpected '!'".to_string()));
                    }
                }
                '<' => {
                    if self.peek() == '=' {
                        self.pos += 2;
                        self.tokens.push(Token::LessThanOrEqual);
                    } else {
                        self.pos += 1;
                        self.tokens.push(Token::LessThan);
                    }
                }
                '>' => {
                    if self.peek() == '=' {
                        self.pos += 2;
                        self.tokens.push(Token::GreaterThanOrEqual);
                    } else {
                        self.pos += 1;
                        self.tokens.push(Token::GreaterThan);
                    }
                }
                '"' | '\'' => {
                    self.tokenize_string();
                }
                c if c.is_ascii_digit() => {
                    self.tokenize_number();
                }
                c if c.is_alphabetic() || c == '_' => {
                    self.tokenize_identifier();
                }
                _ => {
                    self.pos += 1;
                    self.tokens.push(Token::Error(format!("unexpected character: {}", c)));
                }
            }
        }
        
        self.tokens.push(Token::Eof);
        self.tokens.clone()
    }
    
    fn current_char(&self) -> char {
        self.input.as_bytes()[self.pos] as char
    }
    
    fn peek(&self) -> char {
        if self.pos + 1 >= self.input.len() {
            return '\0';
        }
        self.input.as_bytes()[self.pos + 1] as char
    }
    
    fn skip_whitespace(&mut self) {
        while self.pos < self.input.len() && self.input.as_bytes()[self.pos].is_ascii_whitespace() {
            self.pos += 1;
        }
    }
    
    fn tokenize_number(&mut self) {
        let start = self.pos;
        let mut has_dot = false;
        let mut has_exp = false;
        
        while self.pos < self.input.len() {
            let c = self.input.as_bytes()[self.pos] as char;
            if c.is_ascii_digit() {
                self.pos += 1;
            } else if c == '.' && !has_dot {
                has_dot = true;
                self.pos += 1;
            } else if (c == 'e' || c == 'E') && !has_exp {
                has_exp = true;
                self.pos += 1;
                if self.pos < self.input.len() && (self.input.as_bytes()[self.pos] as char == '+' || self.input.as_bytes()[self.pos] as char == '-') {
                    self.pos += 1;
                }
            } else {
                break;
            }
        }
        
        let num_str = &self.input[start..self.pos];
        match num_str.parse::<f64>() {
            Ok(num) => self.tokens.push(Token::Number(num)),
            Err(_) => self.tokens.push(Token::Error(format!("invalid number: {}", num_str))),
        }
    }
    
    fn tokenize_string(&mut self) {
        let quote = self.current_char();
        self.pos += 1; // skip opening quote
        
        let mut result = String::new();
        while self.pos < self.input.len() {
            let c = self.current_char();
            if c == quote {
                self.pos += 1; // skip closing quote
                self.tokens.push(Token::String(result));
                return;
            }
            if c == '\\' {
                self.pos += 1;
                if self.pos < self.input.len() {
                    let escaped = self.current_char();
                    match escaped {
                        'n' => result.push('\n'),
                        't' => result.push('\t'),
                        'r' => result.push('\r'),
                        '\\' => result.push('\\'),
                        '"' => result.push('"'),
                        '\'' => result.push('\''),
                        _ => {
                            result.push('\\');
                            result.push(escaped);
                        }
                    }
                    self.pos += 1;
                }
            } else {
                result.push(c);
                self.pos += 1;
            }
        }
        
        self.tokens.push(Token::Error("unterminated string".to_string()));
    }
    
    fn tokenize_identifier(&mut self) {
        let start = self.pos;
        while self.pos < self.input.len() {
            let c = self.input.as_bytes()[self.pos] as char;
            if c.is_alphanumeric() || c == '_' {
                self.pos += 1;
            } else {
                break;
            }
        }
        
        let ident = &self.input[start..self.pos];
        
        // Check for keywords
        let token = match ident {
            "by" => Token::By,
            "without" => Token::Without,
            "on" => Token::On,
            "ignoring" => Token::Ignoring,
            "group_left" => Token::GroupLeft,
            "group_right" => Token::GroupRight,
            "bool" => Token::Bool,
            "offset" => Token::Offset,
            _ => Token::Identifier(ident.to_string()),
        };
        
        self.tokens.push(token);
    }
}

/// Convenience function to tokenize
pub fn tokenize(input: &str) -> Vec<Token> {
    let mut lexer = Lexer::new(input.to_string());
    lexer.tokenize()
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_tokenize_number() {
        let tokens = tokenize("42");
        assert_eq!(tokens[0], Token::Number(42.0));
    }
    
    #[test]
    fn test_tokenize_float() {
        let tokens = tokenize("3.14");
        assert_eq!(tokens[0], Token::Number(3.14));
    }
    
    #[test]
    fn test_tokenize_string() {
        let tokens = tokenize("\"hello\"");
        assert_eq!(tokens[0], Token::String("hello".to_string()));
    }
    
    #[test]
    fn test_tokenize_identifier() {
        let tokens = tokenize("cpu_usage");
        assert_eq!(tokens[0], Token::Identifier("cpu_usage".to_string()));
    }
    
    #[test]
    fn test_tokenize_operators() {
        let tokens = tokenize("+ - * / % ^");
        assert_eq!(tokens[0], Token::Plus);
        assert_eq!(tokens[1], Token::Minus);
        assert_eq!(tokens[2], Token::Star);
        assert_eq!(tokens[3], Token::Slash);
        assert_eq!(tokens[4], Token::Percent);
        assert_eq!(tokens[5], Token::Caret);
    }
    
    #[test]
    fn test_tokenize_comparison() {
        let tokens = tokenize("= != < > <= >=");
        assert_eq!(tokens[0], Token::Equal);
        assert_eq!(tokens[1], Token::NotEqual);
        assert_eq!(tokens[2], Token::LessThan);
        assert_eq!(tokens[3], Token::GreaterThan);
        assert_eq!(tokens[4], Token::LessThanOrEqual);
        assert_eq!(tokens[5], Token::GreaterThanOrEqual);
    }
    
    #[test]
    fn test_tokenize_regex() {
        let tokens = tokenize("=~ !~");
        assert_eq!(tokens[0], Token::RegexMatch);
        assert_eq!(tokens[1], Token::RegexNotMatch);
    }
    
    #[test]
    fn test_tokenize_delimiters() {
        let tokens = tokenize("(){}[]");
        assert_eq!(tokens[0], Token::LeftParen);
        assert_eq!(tokens[1], Token::RightParen);
        assert_eq!(tokens[2], Token::LeftBrace);
        assert_eq!(tokens[3], Token::RightBrace);
        assert_eq!(tokens[4], Token::LeftBracket);
        assert_eq!(tokens[5], Token::RightBracket);
    }
    
    #[test]
    fn test_tokenize_keywords() {
        let tokens = tokenize("by without on ignoring");
        assert_eq!(tokens[0], Token::By);
        assert_eq!(tokens[1], Token::Without);
        assert_eq!(tokens[2], Token::On);
        assert_eq!(tokens[3], Token::Ignoring);
    }
    
    #[test]
    fn test_tokenize_complex() {
        let tokens = tokenize("rate(cpu_usage{host=\"server1\"}[5m])");
        assert_eq!(tokens[0], Token::Identifier("rate".to_string()));
        assert_eq!(tokens[1], Token::LeftParen);
        assert_eq!(tokens[2], Token::Identifier("cpu_usage".to_string()));
        assert_eq!(tokens[3], Token::LeftBrace);
        assert_eq!(tokens[4], Token::Identifier("host".to_string()));
        assert_eq!(tokens[5], Token::Equal);
        assert_eq!(tokens[6], Token::String("server1".to_string()));
        assert_eq!(tokens[7], Token::RightBrace);
        assert_eq!(tokens[8], Token::LeftBracket);
        assert_eq!(tokens[9], Token::Number(5.0));
        assert_eq!(tokens[10], Token::Identifier("m".to_string()));
    }
    
    #[test]
    fn test_tokenize_empty() {
        let tokens = tokenize("");
        assert_eq!(tokens.len(), 1);
        assert_eq!(tokens[0], Token::Eof);
    }
    
    #[test]
    fn test_tokenize_escape_sequences() {
        let tokens = tokenize("\"hello\\nworld\"");
        assert_eq!(tokens[0], Token::String("hello\nworld".to_string()));
    }
}