#!/bin/bash

# Test 1: Método no permitido (GET en lugar de POST)
echo "Test 1: Método no permitido (GET)"
curl -v http://localhost:8080/api/execute

echo -e "\n\n"

# Test 2: Content-Type incorrecto
echo "Test 2: Content-Type incorrecto"
curl -v -X POST -H "Content-Type: text/plain" http://localhost:8080/api/execute -d "código inválido"

echo -e "\n\n"

# Test 3: JSON inválido
echo "Test 3: JSON inválido"
curl -v -X POST -H "Content-Type: application/json" http://localhost:8080/api/execute -d "{código inválido}"

echo -e "\n\n"

# Test 4: Código vacío
echo "Test 4: Código vacío"
curl -v -X POST -H "Content-Type: application/json" http://localhost:8080/api/execute -d '{"code":""}'

echo -e "\n\n"

# Test 5: Código con import prohibido
echo "Test 5: Código con import prohibido"
curl -v -X POST -H "Content-Type: application/json" http://localhost:8080/api/execute -d '{"code":"package main\nimport \"os/exec\"\nfunc main() {\n\tcmd := exec.Command(\"ls\")\n\tcmd.Run()\n}"}'

echo -e "\n\n"

# Test 6: Código válido
echo "Test 6: Código válido"
curl -v -X POST -H "Content-Type: application/json" http://localhost:8080/api/execute -d '{"code":"package main\nimport \"fmt\"\nfunc main() {\n\tfmt.Println(\"Hello, World!\")\n}"}'
