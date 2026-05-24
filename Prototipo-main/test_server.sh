#!/bin/bash

echo "🧪 Testando configuração do servidor..."
echo ""

# Verificar se os arquivos HTML existem
echo "📄 Verificando arquivos HTML..."
for file in lista.html calendario.html formulariotreinamento.html detalhesdotreinamento.html detalhesdaloja.html engajamento.html; do
    if [ -f "src/imports/$file" ]; then
        echo "  ✅ $file"
    else
        echo "  ❌ $file (NÃO ENCONTRADO)"
    fi
done

echo ""
echo "🖼️  Verificando assets..."
if [ -f "assets/logo_2024_(1)-nvcRYmem.png" ]; then
    echo "  ✅ logo_2024_(1)-nvcRYmem.png"
else
    echo "  ❌ logo_2024_(1)-nvcRYmem.png (NÃO ENCONTRADO)"
fi

echo ""
echo "📦 Verificando main.go..."
if [ -f "main.go" ]; then
    echo "  ✅ main.go"
else
    echo "  ❌ main.go (NÃO ENCONTRADO)"
fi

echo ""
echo "✅ Verificação completa!"
echo ""
echo "Para iniciar o servidor, execute:"
echo "  go run main.go"
