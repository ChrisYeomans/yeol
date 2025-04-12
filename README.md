## Yeol Lang

#### Grammer
```text
block = { instr[] }
term = <input> | variable | literal
expression = term | term + term | term - term | term / term | term * term | term % term
rel = term < term | term > term | term <= term | term >= term | term == term | term != term
instr = (let type)? variable = expression | <if> rel <then> block (<else> block)? | <print> term
method = method methodName(param: type): returnType block
```

