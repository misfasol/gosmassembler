package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Token_tipo string

const (
    label Token_tipo = "label"
    ref Token_tipo = "ref"
    inst Token_tipo = "inst"
    reg Token_tipo = "reg"
    lit Token_tipo = "lit"
)

type Reg int

const (
    rcp Reg = 0x01
    rax Reg = 0x02
    rbx Reg = 0x03
    rcx Reg = 0x04
    rdx Reg = 0x05
)

type Inst int

const (
    addrb Inst = 0x00
    addrr Inst = 0x01
    incr Inst = 0x10
    decr Inst = 0x20
    subrb Inst = 0x30
    subrr Inst = 0x31
    movrb Inst = 0x40
    movrr Inst = 0x41
    jmpb Inst = 0x50
    cmprb Inst = 0x60
    cmprr Inst = 0x61
    jzb Inst = 0x70
)

type Token struct {
    tipo Token_tipo
    texto string
}

func comp(r rune) bool {
    return r == ' ' || r == '\n' || r == '\r' || r == ','
}

func separarPalavras(s string) []string {
    palavras := strings.FieldsFunc(s, comp)
    for i, palavra := range palavras {
        palavras[i] = strings.TrimSpace(palavra)
    }
    return palavras
}

func conseguirLabels(palavras []string) map[string]int {
    label_map := make(map[string]int)
    for i, linha := range palavras {
        if strings.HasSuffix(linha, ":") {
            label_map[strings.TrimSuffix(linha, ":")] = i - len(label_map)
        }
    }
    return label_map
}

func checarLabel(s string) bool {
    return strings.HasSuffix(s, ":")
}

func checarNum(s string) bool {
    for _, r := range s {
        if !unicode.IsDigit(rune(r)) {
            return false
        }
    }
    return true
}

func checarReg(s string) bool {
    regs := [5]string{"rcp", "rax", "rbx", "rcx", "rdx"}
    for _, reg := range regs {
        if s == reg {
            return true
        }
    }
    return false
}

func checarInst(s string) bool {
    insts := [8]string{"add", "inc", "dec", "sub", "mov", "jmp", "cmp", "jz"}
    for _, inst := range insts {
        if s == inst {
            return true
        }
    }
    return false
}

func checarRef(s string, m map[string]int) bool {
    _, ok := m[s]
    return ok
}

func tokenizar(palavras []string, label_map map[string]int) []Token {
    var token_list []Token
    for i, palavra := range palavras {
        if palavra == "" {
            continue
        }
        if checarLabel(palavra) {
            token_list = append(token_list, Token{tipo: label, texto: strings.TrimSuffix(palavra, ":")})
        } else if checarNum(palavra) {
            token_list = append(token_list, Token{tipo: lit, texto: palavra})
        } else if checarReg(palavra) {
            token_list = append(token_list, Token{tipo: reg, texto: palavra})
        } else if checarInst(palavra) {
            token_list = append(token_list, Token{tipo: inst, texto: palavra})
        } else if checarRef(palavra, label_map) {
            valor, _ := label_map[palavra]
            token_list = append(token_list, Token{tipo: ref, texto: fmt.Sprint(valor)})
        } else {
            fmt.Printf("Texto não reconhecido no token: \"%v\", token número: %v", palavra, i)
            os.Exit(1)
        }
    }
    return token_list
}

func escolherInstEspecifica(t Token_tipo, i1, i2 Inst) Inst {
    if t == lit {
        return i1
    } else {
        return i2
    }
}

func escolherInst(s string, t Token_tipo) Inst {
    switch s {
    case "add":
        return escolherInstEspecifica(t, addrb, addrr)
    case "inc":
        return incr
    case "dec":
        return decr
    case "sub":
        return escolherInstEspecifica(t, subrb, subrr)
    case "mov":
        return escolherInstEspecifica(t, movrb, movrr)
    case "jmp":
        return jmpb
    case "cmp":
        return escolherInstEspecifica(t, cmprb, cmprr)
    case "jz":
        return jzb
    default:
        fmt.Printf("Instrução não reconhecida: %v\n", s)
        os.Exit(1) // o código deveria nunca chegar aqui
        return 0
    }
}

func escolherReg(s string) Reg {
    switch s {
    case "rcp":
        return rcp
    case "rax":
        return rax
    case "rbx":
        return rbx
    case "rcx":
        return rcx
    case "rdx":
        return rdx
    default: // em teoria o código nunca deveria chegar aqui
        return 0
    }
}

func criarBinario(nome string, token_list []Token) {
    arq_bin, err := os.Create(nome)
    if err != nil {
        fmt.Println("Erro criando arquivo binário")
        os.Exit(1)
    }
    
    var bit []byte

    for i, v := range token_list {
        switch v.tipo {
        case label:
            // nada a fazer aqui
        case inst:
            bit = append(bit, byte(escolherInst(v.texto, token_list[i+2].tipo)))
        case reg:
            bit = append(bit, byte(escolherReg(v.texto)))
        case lit:
            valor, _ := strconv.Atoi(v.texto)
            bit = append(bit, byte(valor))
        case ref:
            valor, _ := strconv.Atoi(v.texto)
            bit = append(bit, byte(valor))
        }
    }

    arq_bin.Write(bit)
}

func main() {
    arr, err := os.ReadFile(os.Args[1])
    if err != nil {
        fmt.Println("Erro abrindo arquivo")
        os.Exit(0)
    }

    palavras := separarPalavras(string(arr))

    label_map := conseguirLabels(palavras)

    token_list := tokenizar(palavras, label_map)

    criarBinario(strings.TrimSuffix(os.Args[1], ".asm") + ".bin", token_list)
}
