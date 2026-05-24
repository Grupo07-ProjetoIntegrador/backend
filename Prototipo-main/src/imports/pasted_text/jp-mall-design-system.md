🎨 Design System: JP Mall Corporativo
1. Visão Geral
Este Design System foi criado para garantir que a interface do sistema web corporativo JP Mall seja limpa, moderna, altamente funcional e perfeitamente alinhada à identidade visual (branding) da marca (Grupo Flamboyant). O foco principal é fornecer uma ferramenta de trabalho robusta, acessível e sem distrações desnecessárias.

2. Paleta de Cores (Color Tokens)
As cores foram escolhidas para transmitir seriedade, elegância corporativa e contrastes bem definidos para usabilidade diária.
🔴 Cores da Marca (Brand Colors)
Vermelho Vinho Escuro (Primária): #8B1A1A
Uso: Menu lateral (Navegação principal), Cabeçalhos, Ícones de destaque neutro.
Vermelho Vibrante (Ação/Destaque): #D93030
Uso: Botões primários (ex: "Salvar", "Entrar"), Links de ação, Alertas críticos e Hovers interativos.
Bege / Areia Claro (Fundo): #F7F4EF
Uso: Cor de fundo padrão de todas as telas (background-color da aplicação).
Bege Escuro / Dourado (Apoio): #C8A882
Uso: Detalhes gráficos, bordas secundárias, gráficos de dados.
⚪ Cores Neutras (Neutral Colors)
Branco Puro: #FFFFFF
Uso: Fundo exclusivo para Cards, Tabelas, Formulários e Modais.
Chumbo / Cinza Escuro: #1F2937 (Tailwind gray-800 / gray-900)
Uso: Títulos (H1, H2, H3), textos principais de dados.
Cinza Médio: #4B5563 a #6B7280 (Tailwind gray-500 a gray-600)
Uso: Textos descritivos, subtítulos, labels de formulário, placeholders.
Cinza Claro (Bordas): #E5E7EB (Tailwind gray-200)
Uso: Bordas de cards, divisórias e campos de input.
🚦 Cores Semânticas (Feedback Colors)
Sucesso / Aprovado: Verde (#10B981 / Tailwind green-500 a green-700)
Atenção / Em Análise: Amarelo/Laranja (#F59E0B / Tailwind yellow-500 ou orange-500)
Erro / Recusado: Vermelho Vibrante (#D93030)
Informativo: Azul (#3B82F6 / Tailwind blue-600)

3. Tipografia (Typography)
A tipografia prioriza a legibilidade para longas jornadas de trabalho em telas de desktop.
Fonte Principal: Inter (sans-serif)
Alternativa/Fallback: Roboto, system-ui, sans-serif
Tamanho Base (Body): 16px (Regular / 400)
Escala Tipográfica:
H1 (Títulos de Página): 24px (Bold / 700), cor Chumbo (gray-800)
H2 (Títulos de Seção/Cards): 18px a 20px (Bold / 700), cor Chumbo
H3 (Subtítulos/Labels fortes): 14px (Bold / 700) ou 12px (Uppercase, Bold, tracking-wider)
Body (Texto Corrido): 14px a 16px (Regular / 400), cor Cinza Médio
Small (Apoio/Dicas): 12px (Regular / 400), cor Cinza Médio

4. Espaçamento e Grid (Spacing Rules)
O layout obedece à regra estrita de múltiplos de 8px (8pt grid system) para garantir ritmo vertical e horizontal perfeito.
Micro: 4px (gap-1, mb-1 - Entre ícone e texto)
Pequeno: 8px (gap-2, p-2 - Margens internas de botões e inputs)
Médio: 16px (p-4, mb-4 - Espaçamento padrão entre componentes)
Grande: 24px (p-6, gap-6 - Paddings internos de Cards)
Extra Grande: 32px (p-8, mb-8 - Espaçamento entre seções da página)
Bordas Arredondadas (Border Radius):
Padrão: 8px (rounded-lg) para Cards, Inputs, Botões e Modais.

5. Componentes UI (UI Components)
🔘 Botões (Buttons)
Botão Primário: Fundo #D93030, Texto Branco, sem bordas. Efeito hover para um tom levemente mais escuro (#b92828).
Botão Secundário / Cancelar: Fundo Branco (ou Cinza muito claro gray-50), Borda Cinza Clara (gray-200), Texto Cinza Escuro.
Tamanho (Acessibilidade): Todos os botões clicáveis devem ter altura mínima de 44px (min-h-[44px]).
📝 Inputs e Formulários (Forms)
Fundo: Cinza muito claro (gray-50) quando inativo, mudando para Branco (bg-white) quando em foco (focus).
Bordas: Finas e sutis (border-gray-200).
Foco (Focus State): Borda colorida em Vermelho Vibrante (focus:ring-2 focus:ring-[#D93030] focus:border-transparent) para clareza visual de onde o cursor está.
Labels: Posicionadas na parte superior externa do campo, em texto menor (14px), peso médio (font-medium) e cor cinza (text-gray-700).
🗂️ Cards e Contêineres (Containers)
Fundo: Sempre Branco.
Borda: Muito sutil (border-gray-200).
Sombra (Shadow): Drop shadow leve (shadow-sm) para criar profundidade e fazer o card saltar ligeiramente do fundo Bege (#F7F4EF) da tela.
Cantos: Arredondados em 8px (rounded-lg).
🏷️ Tags de Status (Badges)
Utilizadas em tabelas e cabeçalhos para identificação rápida. Têm fundo com 10-20% de opacidade e texto na cor sólida.
Verde (Pago/Aprovado): Fundo green-50, Texto green-700, Borda green-200.
Amarelo/Laranja (Aguardando/Em Análise): Fundo orange-50, Texto orange-700, Borda orange-200.
Vermelho (Alta Gravidade): Fundo red-50, Texto #D93030, Borda red-200.
Azul (Aberto/Novo): Fundo blue-50, Texto blue-700, Borda blue-200.

6. Diretrizes de Acessibilidade (WCAG)
Contraste: A regra de texto escuro sobre fundos claros, e texto branco sobre o vermelho escuro (#8B1A1A) ou vibrante (#D93030), garante taxa de contraste de aprovação WCAG AA/AAA.
Área de Toque (Touch Targets): Inputs, selects e botões obedecem à regra de altura mínima de 44px (ou superior) para facilitar cliques em touchscreens (foco em uso de tablets/dispositivos móveis pelos lojistas).
Estados de Foco: Todos os elementos interativos possuem um estado de foco (outline/ring) claro e visível em vermelho vibrante.
Hierarquia Visual: Uso claro da semântica HTML e hierarquia de tamanhos de fonte em negrito para guiar o olho rapidamente pelos painéis e dados complexos.

