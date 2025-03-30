## Yeol Lang

#### Grammer
```text
term = <input> | variable | literal
expression = term | term + term | term - term | term / term | term * term | term % term
rel = term < term | term > term | term <= term | term >= term | term == term | term != term
instr = variable = expression | <if> rel <then> instr | <print> term
```