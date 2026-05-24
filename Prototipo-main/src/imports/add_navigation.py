#!/usr/bin/env python3
"""
Script para adicionar navegação onclick nos arquivos HTML estáticos
"""

import re

def add_onclick_lista():
    """Adiciona onclick em lista.html"""
    with open('lista.html', 'r', encoding='utf-8') as f:
        content = f.read()

    # 1. Botão "Novo Treinamento"
    content = re.sub(
        r'(<button[^>]*)(style="background-color: rgb\(217, 48, 48\);">[\s\S]*?Novo Treinamento</button>)',
        r'\1 onclick="window.location.href=\'/treinamentos/novo\'" \2',
        content,
        count=1
    )

    # 2. Aba "Engajamento das Lojas"
    content = re.sub(
        r'(<button[^>]*data-radix-collection-item="">)(Engajamento das Lojas</button>)',
        r'<button onclick="window.location.href=\'/engajamento\'" \1\2',
        content
    )

    # 3. Botão "Calendário" (toggle)
    content = re.sub(
        r'(<button class="px-3 py-1\.5 text-sm font-medium rounded-md flex items-center gap-2 transition-colors[^>]*>[\s\S]*?</svg>)(Calendário</button>)',
        r'<button onclick="window.location.href=\'/treinamentos/calendario\'" \1\2',
        content,
        count=1
    )

    # 4. Linhas da tabela (tr) - adicionar onclick
    content = re.sub(
        r'(<tr[^>]*class="[^"]*border-b border-gray-100 hover:bg-gray-50[^"]*"[^>]*)(>)',
        r'\1 onclick="window.location.href=\'/treinamentos/detalhes\'" style="cursor: pointer"\2',
        content
    )

    # 5. Botão editar (prevenir propagação)
    content = re.sub(
        r'(<button[^>]*title="Editar treinamento")',
        r'\1 onclick="event.stopPropagation(); window.location.href=\'/treinamentos/editar\'"',
        content
    )

    # 6. Botão excluir (prevenir propagação - não funcional mas mantém)
    content = re.sub(
        r'(<button[^>]*title="Excluir treinamento")',
        r'\1 onclick="event.stopPropagation(); alert(\'Funcionalidade indisponível em HTML estático\')"',
        content
    )

    with open('lista.html', 'w', encoding='utf-8') as f:
        f.write(content)
    print("✓ lista.html atualizado")


def add_onclick_calendario():
    """Adiciona onclick em calendario.html"""
    with open('calendario.html', 'r', encoding='utf-8') as f:
        content = f.read()

    # 1. Botão "Novo Treinamento"
    content = re.sub(
        r'(<button[^>]*)(style="background-color: rgb\(217, 48, 48\);">[\s\S]*?Novo Treinamento</button>)',
        r'\1 onclick="window.location.href=\'/treinamentos/novo\'" \2',
        content,
        count=1
    )

    # 2. Botão "Lista" (toggle)
    content = re.sub(
        r'(<button class="px-3 py-1\.5 text-sm font-medium rounded-md flex items-center gap-2 transition-colors[^>]*>[\s\S]*?</svg>)(Lista</button>)',
        lambda m: f'<button onclick="window.location.href=\'/treinamentos/lista\'" {m.group(1)}{m.group(2)}',
        content,
        count=1
    )

    # 3. Aba "Engajamento das Lojas"
    content = re.sub(
        r'(<button[^>]*data-radix-collection-item="">)(Engajamento das Lojas</button>)',
        r'<button onclick="window.location.href=\'/engajamento\'" \1\2',
        content
    )

    # 4. Cards de treinamento no calendário (com ícones de status)
    content = re.sub(
        r'(<button[^>]*class="text-left px-2 py-1 text-xs rounded-md border[^>]*)(>)',
        r'\1 onclick="window.location.href=\'/treinamentos/detalhes\'"\2',
        content
    )

    with open('calendario.html', 'w', encoding='utf-8') as f:
        f.write(content)
    print("✓ calendario.html atualizado")


def add_onclick_formulario():
    """Adiciona onclick em formulariotreinamento.html"""
    with open('formulariotreinamento.html', 'r', encoding='utf-8') as f:
        content = f.read()

    # 1. Botão "Voltar"
    content = re.sub(
        r'(<button[^>]*>[\s\S]*?</svg>)(Voltar</button>)',
        r'<button onclick="window.location.href=\'/treinamentos/lista\'" \1\2',
        content,
        count=1
    )

    # 2. Botão "Cancelar"
    content = re.sub(
        r'(<button[^>]*>)(Cancelar</button>)',
        r'<button onclick="window.location.href=\'/treinamentos/lista\'" \1\2',
        content
    )

    # 3. Botão "Salvar" / "Criar" (não funcional)
    content = re.sub(
        r'(<button[^>]*data-slot="button"[^>]*style="background-color: rgb\(217, 48, 48\);[^>]*>)(Salvar|Criar Treinamento)',
        r'<button onclick="alert(\'Formulário não funcional em HTML estático\'); window.location.href=\'/treinamentos/lista\'" \1\2',
        content
    )

    with open('formulariotreinamento.html', 'w', encoding='utf-8') as f:
        f.write(content)
    print("✓ formulariotreinamento.html atualizado")


def add_onclick_detalhes_treinamento():
    """Adiciona onclick em detalhesdotreinamento.html"""
    with open('detalhesdotreinamento.html', 'r', encoding='utf-8') as f:
        content = f.read()

    # 1. Botão "Voltar"
    content = re.sub(
        r'(<button[^>]*>[\s\S]*?</svg>)(Voltar</button>)',
        r'<button onclick="window.location.href=\'/treinamentos/lista\'" \1\2',
        content,
        count=1
    )

    # 2. Botão "Editar"
    content = re.sub(
        r'(<button[^>]*>[\s\S]*?</svg>)(Editar</button>)',
        r'<button onclick="window.location.href=\'/treinamentos/editar\'" \1\2',
        content
    )

    # 3. Botão "Excluir" (não funcional)
    content = re.sub(
        r'(<button[^>]*>[\s\S]*?</svg>)(Excluir</button>)',
        r'<button onclick="alert(\'Funcionalidade indisponível em HTML estático\')" \1\2',
        content
    )

    with open('detalhesdotreinamento.html', 'w', encoding='utf-8') as f:
        f.write(content)
    print("✓ detalhesdotreinamento.html atualizado")


def add_onclick_engajamento():
    """Adiciona onclick em engajamento.html"""
    with open('engajamento.html', 'r', encoding='utf-8') as f:
        content = f.read()

    # 1. Aba "Agenda e Lista"
    content = re.sub(
        r'(<button[^>]*data-state="inactive"[^>]*data-radix-collection-item="">)(Agenda[\s\S]*?e Lista</button>)',
        r'<button onclick="window.location.href=\'/treinamentos/lista\'" \1\2',
        content,
        count=1
    )

    # 2. Botão "Ver Detalhes" (lojas)
    content = re.sub(
        r'(<button[^>]*data-slot="button"[^>]*>)(Ver Detalhes</button>)',
        r'<button onclick="event.stopPropagation(); window.location.href=\'/lojas/detalhes\'" \1\2',
        content
    )

    # 3. Linhas da tabela de lojas (tr)
    content = re.sub(
        r'(<tr[^>]*class="cursor-pointer hover:bg-gray-50[^>]*)(>)',
        r'\1 onclick="window.location.href=\'/lojas/detalhes\'"\2',
        content
    )

    with open('engajamento.html', 'w', encoding='utf-8') as f:
        f.write(content)
    print("✓ engajamento.html atualizado")


def add_onclick_detalhes_loja():
    """Adiciona onclick em detalhesdaloja.html"""
    with open('detalhesdaloja.html', 'r', encoding='utf-8') as f:
        content = f.read()

    # 1. Botão "Voltar"
    content = re.sub(
        r'(<button[^>]*>[\s\S]*?</svg>)(Voltar</button>)',
        r'<button onclick="window.location.href=\'/engajamento\'" \1\2',
        content,
        count=1
    )

    # 2. Linhas de treinamentos clicáveis
    content = re.sub(
        r'(<tr[^>]*class="cursor-pointer[^>]*)(>)',
        r'\1 onclick="window.location.href=\'/treinamentos/detalhes\'"\2',
        content
    )

    with open('detalhesdaloja.html', 'w', encoding='utf-8') as f:
        f.write(content)
    print("✓ detalhesdaloja.html atualizado")


if __name__ == '__main__':
    print("🚀 Adicionando navegação onclick nos arquivos HTML...\n")

    try:
        add_onclick_lista()
        add_onclick_calendario()
        add_onclick_formulario()
        add_onclick_detalhes_treinamento()
        add_onclick_engajamento()
        add_onclick_detalhes_loja()

        print("\n✅ Todos os arquivos foram atualizados com sucesso!")
        print("\n📋 Próximos passos:")
        print("   1. Configure o main.go com as rotas")
        print("   2. Teste a navegação entre as páginas")

    except Exception as e:
        print(f"\n❌ Erro: {e}")
