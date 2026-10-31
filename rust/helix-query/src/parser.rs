//! Helix Query Engine - Parser Module
//! 
//! This module provides a complete parser for the Helix query language.
//! Supports instant vectors, range vectors, aggregations, and binary operations.

use std::collections::BTreeMap;
use std::fmt;

use super::lexer::{Token, tokenize};

/// AST node types
#[derive(Debug, Clone, PartialEq)]
pub enum Expr {
    // Literals
    Number(f64),
    String(String),
    
    // Identifiers and selectors
    Identifier(String),
    VectorSelector {
        name: String,
        matchers: Vec<Matcher>,
    },
    MatrixSelector {
        name: String,
        matchers: Vec<Matcher>,
        range: Duration,
        offset: Option<Duration>,
    },
    
    // Aggregations
    Aggregation {
        op: AggregationOp,
        expr: Box<Expr>,
        grouping: Vec<String>,
        without: bool,
    },
    
    // Function calls
    FunctionCall {
        name: String,
        args: Vec<Expr>,
    },
    
    // Binary operations
    BinaryOp {
        op: BinaryOperator,
        left: Box<Expr>,
        right: Box<Expr>,
        bool_modifier: bool,
    },
    
    // Unary operations
    UnaryOp {
        op: UnaryOperator,
        expr: Box<Expr>,
    },
    
    // Parentheses
    Paren(Box<Expr>),
    
    // Subqueries
    Subquery {
        expr: Box<Expr>,
        range: Duration,
        step: Option<Duration>,
        offset: Option<Duration>,
    },
}

/// Matcher for label filtering
#[derive(Debug, Clone, PartialEq)]
pub struct Matcher {
    pub name: String,
    pub op: MatchOp,
    pub value: String,
}

/// Match operators
#[derive(Debug, Clone, Copy, PartialEq)]
pub enum MatchOp {
    Equal,
    NotEqual,
    RegexMatch,
    RegexNotMatch,
}

/// Aggregation operators
#[derive(Debug, Clone, Copy, PartialEq)]
pub enum AggregationOp {
    Sum,
    Min,
    Max,
    Avg,
    Count,
    Stddev,
    Stdvar,
    TopK,
    BottomK,
    Quantile,
    Group,
    LimitK,
    LimitRatio,
}

/// Binary operators
#[derive(Debug, Clone, Copy, PartialEq)]
pub enum BinaryOperator {
    Add,
    Sub,
    Mul,
    Div,
    Mod,
    Pow,
    Equal,
    NotEqual,
    GreaterThan,
    LessThan,
    GreaterThanOrEqual,
    LessThanOrEqual,
    And,
    Or,
    Unless,
    Atan2,
}

/// Unary operators
#[derive(Debug, Clone, Copy, PartialEq)]
pub enum UnaryOperator {
    Plus,
    Minus,
}

/// Duration representation
#[derive(Debug, Clone, PartialEq)]
pub struct Duration {
    pub value: f64,
    pub unit: DurationUnit,
}

#[derive(Debug, Clone, Copy, PartialEq)]
pub enum DurationUnit {
    Milliseconds,
    Seconds,
    Minutes,
    Hours,
    Days,
    Weeks,
    Years,
}

impl Duration {
    pub fn to_seconds(&self) -> f64 {
        match self.unit {
            DurationUnit::Milliseconds => self.value / 1000.0,
            DurationUnit::Seconds => self.value,
            DurationUnit::Minutes => self.value * 60.0,
            DurationUnit::Hours => self.value * 3600.0,
            DurationUnit::Days => self.value * 86400.0,
            DurationUnit::Weeks => self.value * 604800.0,
            DurationUnit::Years => self.value * 31536000.0,
        }
    }
}

/// Parser for the query language
pub struct QueryParser {
    tokens: Vec<Token>,
    pos: usize,
}

impl QueryParser {
    pub fn new(input: &str) -> Self {
        let tokens = tokenize(input);
        QueryParser { tokens, pos: 0 }
    }
    
    pub fn parse(&mut self) -> Result<Expr, ParseError> {
        let expr = self.parse_expr()?;
        
        if self.current() != &Token::Eof {
            return Err(ParseError::UnexpectedToken(format!("{:?}", self.current())));
        }
        
        Ok(expr)
    }
    
    fn parse_expr(&mut self) -> Result<Expr, ParseError> {
        self.parse_binary_op(0)
    }
    
    fn parse_binary_op(&mut self, min_prec: u8) -> Result<Expr, ParseError> {
        let mut left = self.parse_unary_op()?;
        
        while self.pos < self.tokens.len() {
            let op = match self.current() {
                Token::Plus => BinaryOperator::Add,
                Token::Minus => BinaryOperator::Sub,
                Token::Star => BinaryOperator::Mul,
                Token::Slash => BinaryOperator::Div,
                Token::Percent => BinaryOperator::Mod,
                Token::Caret => BinaryOperator::Pow,
                Token::Equal => BinaryOperator::Equal,
                Token::NotEqual => BinaryOperator::NotEqual,
                Token::LessThan => BinaryOperator::LessThan,
                Token::GreaterThan => BinaryOperator::GreaterThan,
                Token::LessThanOrEqual => BinaryOperator::LessThanOrEqual,
                Token::GreaterThanOrEqual => BinaryOperator::GreaterThanOrEqual,
                Token::Identifier(s) if s == "and" => BinaryOperator::And,
                Token::Identifier(s) if s == "or" => BinaryOperator::Or,
                Token::Identifier(s) if s == "unless" => BinaryOperator::Unless,
                Token::Identifier(s) if s == "atan2" => BinaryOperator::Atan2,
                _ => break,
            };
            
            let prec = self.binary_precedence(op);
            if prec < min_prec {
                break;
            }
            
            self.advance();
            let bool_modifier = self.consume_bool();
            let right = self.parse_binary_op(prec + 1)?;
            
            left = Expr::BinaryOp {
                op,
                left: Box::new(left),
                right: Box::new(right),
                bool_modifier,
            };
        }
        
        Ok(left)
    }
    
    fn parse_unary_op(&mut self) -> Result<Expr, ParseError> {
        match self.current() {
            Token::Plus => {
                self.advance();
                let expr = self.parse_unary_op()?;
                Ok(Expr::UnaryOp {
                    op: UnaryOperator::Plus,
                    expr: Box::new(expr),
                })
            }
            Token::Minus => {
                self.advance();
                let expr = self.parse_unary_op()?;
                Ok(Expr::UnaryOp {
                    op: UnaryOperator::Minus,
                    expr: Box::new(expr),
                })
            }
            _ => self.parse_primary(),
        }
    }
    
    fn parse_primary(&mut self) -> Result<Expr, ParseError> {
        match self.current().clone() {
            Token::Number(n) => {
                self.advance();
                Ok(Expr::Number(n))
            }
            Token::String(s) => {
                self.advance();
                Ok(Expr::String(s))
            }
            Token::LeftParen => {
                self.advance();
                let expr = self.parse_expr()?;
                self.expect(Token::RightParen)?;
                Ok(Expr::Paren(Box::new(expr)))
            }
            Token::Identifier(name) => {
                self.advance();
                
                // Check for function call
                if self.current() == &Token::LeftParen {
                    return self.parse_function_call(name);
                }
                
                // Check for vector selector
                if self.current() == &Token::LeftBrace {
                    return self.parse_vector_selector(name);
                }
                
                // Check for matrix selector
                if self.current() == &Token::LeftBracket {
                    return self.parse_matrix_selector(name);
                }
                
                Ok(Expr::Identifier(name))
            }
            _ => Err(ParseError::UnexpectedToken(format!("{:?}", self.current()))),
        }
    }
    
    fn parse_function_call(&mut self, name: String) -> Result<Expr, ParseError> {
        self.expect(Token::LeftParen)?;
        
        let mut args = Vec::new();
        if self.current() != &Token::RightParen {
            args.push(self.parse_expr()?);
            while self.current() == &Token::Comma {
                self.advance();
                args.push(self.parse_expr()?);
            }
        }
        
        self.expect(Token::RightParen)?;
        
        // Check for aggregation
        if self.is_aggregation(&name) {
            return self.parse_aggregation(name, args);
        }
        
        Ok(Expr::FunctionCall { name, args })
    }
    
    fn parse_aggregation(&mut self, op: String, args: Vec<Expr>) -> Result<Expr, ParseError> {
        let agg_op = match op.as_str() {
            "sum" => AggregationOp::Sum,
            "min" => AggregationOp::Min,
            "max" => AggregationOp::Max,
            "avg" => AggregationOp::Avg,
            "count" => AggregationOp::Count,
            "stddev" => AggregationOp::Stddev,
            "stdvar" => AggregationOp::Stdvar,
            "topk" => AggregationOp::TopK,
            "bottomk" => AggregationOp::BottomK,
            "quantile" => AggregationOp::Quantile,
            "group" => AggregationOp::Group,
            _ => return Err(ParseError::UnknownFunction(op)),
        };
        
        // Parse grouping
        let mut grouping = Vec::new();
        let mut without = false;
        
        if self.current() == &Token::By {
            self.advance();
            self.expect(Token::LeftParen)?;
            grouping = self.parse_label_list()?;
            self.expect(Token::RightParen)?;
        } else if self.current() == &Token::Without {
            self.advance();
            without = true;
            self.expect(Token::LeftParen)?;
            grouping = self.parse_label_list()?;
            self.expect(Token::RightParen)?;
        }
        
        Ok(Expr::Aggregation {
            op: agg_op,
            expr: Box::new(args.into_iter().next().unwrap_or(Expr::Number(0.0))),
            grouping,
            without,
        })
    }
    
    fn parse_label_list(&mut self) -> Result<Vec<String>, ParseError> {
        let mut labels = Vec::new();
        
        if self.current() == &Token::RightParen {
            return Ok(labels);
        }
        
        if let Token::Identifier(name) = self.current().clone() {
            labels.push(name);
            self.advance();
        }
        
        while self.current() == &Token::Comma {
            self.advance();
            if let Token::Identifier(name) = self.current().clone() {
                labels.push(name);
                self.advance();
            } else {
                return Err(ParseError::ExpectedIdentifier);
            }
        }
        
        Ok(labels)
    }
    
    fn parse_vector_selector(&mut self, name: String) -> Result<Expr, ParseError> {
        self.expect(Token::LeftBrace)?;
        
        let mut matchers = Vec::new();
        if self.current() != &Token::RightBrace {
            matchers.push(self.parse_matcher()?);
            while self.current() == &Token::Comma {
                self.advance();
                matchers.push(self.parse_matcher()?);
            }
        }
        
        self.expect(Token::RightBrace)?;
        
        Ok(Expr::VectorSelector { name, matchers })
    }
    
    fn parse_matrix_selector(&mut self, name: String) -> Result<Expr, ParseError> {
        self.expect(Token::LeftBracket)?;
        
        let range = self.parse_duration()?;
        self.expect(Token::RightBracket)?;
        
        // Parse optional offset
        let offset = if self.current() == &Token::Offset {
            self.advance();
            Some(self.parse_duration()?)
        } else {
            None
        };
        
        // Parse matchers if present
        let matchers = if self.current() == &Token::LeftBrace {
            self.advance();
            let mut m = Vec::new();
            if self.current() != &Token::RightBrace {
                m.push(self.parse_matcher()?);
                while self.current() == &Token::Comma {
                    self.advance();
                    m.push(self.parse_matcher()?);
                }
            }
            self.expect(Token::RightBrace)?;
            m
        } else {
            Vec::new()
        };
        
        Ok(Expr::MatrixSelector {
            name,
            matchers,
            range,
            offset,
        })
    }
    
    fn parse_matcher(&mut self) -> Result<Matcher, ParseError> {
        let name = if let Token::Identifier(n) = self.current().clone() {
            self.advance();
            n
        } else {
            return Err(ParseError::ExpectedIdentifier);
        };
        
        let op = match self.current() {
            Token::Equal => { self.advance(); MatchOp::Equal }
            Token::NotEqual => { self.advance(); MatchOp::NotEqual }
            Token::RegexMatch => { self.advance(); MatchOp::RegexMatch }
            Token::RegexNotMatch => { self.advance(); MatchOp::RegexNotMatch }
            _ => return Err(ParseError::ExpectedOperator),
        };
        
        let value = if let Token::String(s) = self.current().clone() {
            self.advance();
            s
        } else {
            return Err(ParseError::ExpectedString);
        };
        
        Ok(Matcher { name, op, value })
    }
    
    fn parse_duration(&mut self) -> Result<Duration, ParseError> {
        let mut value = 0.0;
        let mut unit = DurationUnit::Seconds;
        
        if let Token::Number(n) = self.current().clone() {
            value = n;
            self.advance();
        } else {
            return Err(ParseError::ExpectedNumber);
        }
        
        if let Token::Identifier(u) = self.current().clone() {
            unit = match u.as_str() {
                "ms" => DurationUnit::Milliseconds,
                "s" => DurationUnit::Seconds,
                "m" => DurationUnit::Minutes,
                "h" => DurationUnit::Hours,
                "d" => DurationUnit::Days,
                "w" => DurationUnit::Weeks,
                "y" => DurationUnit::Years,
                _ => return Err(ParseError::InvalidDurationUnit(u)),
            };
            self.advance();
        }
        
        Ok(Duration { value, unit })
    }
    
    fn consume_bool(&mut self) -> bool {
        if self.current() == &Token::Bool {
            self.advance();
            true
        } else {
            false
        }
    }
    
    fn is_aggregation(&self, name: &str) -> bool {
        matches!(name, "sum" | "min" | "max" | "avg" | "count" | "stddev" | "stdvar" | "topk" | "bottomk" | "quantile" | "group" | "limitk" | "limit_ratio")
    }
    
    fn binary_precedence(&self, op: BinaryOperator) -> u8 {
        match op {
            BinaryOperator::Pow => 10,
            BinaryOperator::Mul | BinaryOperator::Div | BinaryOperator::Mod | BinaryOperator::Atan2 => 9,
            BinaryOperator::Add | BinaryOperator::Sub => 8,
            BinaryOperator::Equal | BinaryOperator::NotEqual | BinaryOperator::GreaterThan |
            BinaryOperator::LessThan | BinaryOperator::GreaterThanOrEqual | BinaryOperator::LessThanOrEqual => 7,
            BinaryOperator::And => 6,
            BinaryOperator::Or | BinaryOperator::Unless => 5,
            _ => 0,
        }
    }
    
    fn current(&self) -> &Token {
        if self.pos < self.tokens.len() {
            &self.tokens[self.pos]
        } else {
            &Token::Eof
        }
    }
    
    fn advance(&mut self) {
        if self.pos < self.tokens.len() {
            self.pos += 1;
        }
    }
    
    fn expect(&mut self, expected: Token) -> Result<(), ParseError> {
        if self.current() == &expected {
            self.advance();
            Ok(())
        } else {
            Err(ParseError::ExpectedToken(format!("{:?}", expected)))
        }
    }
}

/// Parse errors
#[derive(Debug)]
pub enum ParseError {
    UnexpectedToken(String),
    ExpectedIdentifier,
    ExpectedString,
    ExpectedNumber,
    ExpectedOperator,
    ExpectedToken(String),
    UnknownFunction(String),
    InvalidDurationUnit(String),
}

impl fmt::Display for ParseError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            ParseError::UnexpectedToken(t) => write!(f, "unexpected token: {}", t),
            ParseError::ExpectedIdentifier => write!(f, "expected identifier"),
            ParseError::ExpectedString => write!(f, "expected string"),
            ParseError::ExpectedNumber => write!(f, "expected number"),
            ParseError::ExpectedOperator => write!(f, "expected operator"),
            ParseError::ExpectedToken(t) => write!(f, "expected token: {}", t),
            ParseError::UnknownFunction(f) => write!(f, "unknown function: {}", f),
            ParseError::InvalidDurationUnit(u) => write!(f, "invalid duration unit: {}", u),
        }
    }
}

impl std::error::Error for ParseError {}

/// Convenience function to parse a query
pub fn parse(input: &str) -> Result<Expr, ParseError> {
    let mut parser = QueryParser::new(input);
    parser.parse()
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_parse_number() {
        let expr = parse("42").unwrap();
        assert_eq!(expr, Expr::Number(42.0));
    }
    
    #[test]
    fn test_parse_identifier() {
        let expr = parse("cpu_usage").unwrap();
        assert_eq!(expr, Expr::Identifier("cpu_usage".to_string()));
    }
    
    #[test]
    fn test_parse_vector_selector() {
        let expr = parse("cpu_usage{host=\"server1\"}").unwrap();
        match expr {
            Expr::VectorSelector { name, matchers } => {
                assert_eq!(name, "cpu_usage");
                assert_eq!(matchers.len(), 1);
            }
            _ => panic!("Expected VectorSelector"),
        }
    }
    
    #[test]
    fn test_parse_matrix_selector() {
        let expr = parse("cpu_usage[5m]").unwrap();
        match expr {
            Expr::MatrixSelector { name, range, .. } => {
                assert_eq!(name, "cpu_usage");
                assert_eq!(range.to_seconds(), 300.0);
            }
            _ => panic!("Expected MatrixSelector"),
        }
    }
    
    #[test]
    fn test_parse_function_call() {
        let expr = parse("rate(cpu_usage[5m])").unwrap();
        match expr {
            Expr::FunctionCall { name, args } => {
                assert_eq!(name, "rate");
                assert_eq!(args.len(), 1);
            }
            _ => panic!("Expected FunctionCall"),
        }
    }
    
    #[test]
    fn test_parse_aggregation() {
        let expr = parse("sum by (host) (cpu_usage)").unwrap();
        match expr {
            Expr::Aggregation { op, grouping, without, .. } => {
                assert_eq!(op, AggregationOp::Sum);
                assert_eq!(grouping, vec!["host"]);
                assert!(!without);
            }
            _ => panic!("Expected Aggregation"),
        }
    }
    
    #[test]
    fn test_parse_binary_op() {
        let expr = parse("a + b").unwrap();
        match expr {
            Expr::BinaryOp { op, .. } => {
                assert_eq!(op, BinaryOperator::Add);
            }
            _ => panic!("Expected BinaryOp"),
        }
    }
    
    #[test]
    fn test_parse_unary_op() {
        let expr = parse("-cpu_usage").unwrap();
        match expr {
            Expr::UnaryOp { op, .. } => {
                assert_eq!(op, UnaryOperator::Minus);
            }
            _ => panic!("Expected UnaryOp"),
        }
    }
    
    #[test]
    fn test_parse_paren() {
        let expr = parse("(cpu_usage)").unwrap();
        match expr {
            Expr::Paren(_) => {}
            _ => panic!("Expected Paren"),
        }
    }
    
    #[test]
    fn test_parse_complex() {
        let expr = parse("sum by (host) (rate(cpu_usage{job=\"node\"}[5m]))").unwrap();
        match expr {
            Expr::Aggregation { op, grouping, .. } => {
                assert_eq!(op, AggregationOp::Sum);
                assert_eq!(grouping, vec!["host"]);
            }
            _ => panic!("Expected Aggregation"),
        }
    }
    
    #[test]
    fn test_parse_error() {
        let result = parse("invalid {{");
        assert!(result.is_err());
    }
    
    #[test]
    fn test_duration_to_seconds() {
        let d = Duration { value: 5.0, unit: DurationUnit::Minutes };
        assert_eq!(d.to_seconds(), 300.0);
        
        let d = Duration { value: 1.0, unit: DurationUnit::Hours };
        assert_eq!(d.to_seconds(), 3600.0);
        
        let d = Duration { value: 1000.0, unit: DurationUnit::Milliseconds };
        assert_eq!(d.to_seconds(), 1.0);
    }
}