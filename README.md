# An Architecture Proposal Using Hybrid Blockchain Applied for Supply Chain Tracking

## Submissão de artigos ao IEEE América Latina ( ID: 8343 ) 

A logística da cadeia de abastecimento desempenha um papel crucial na preservação da qualidade e segurança de produtos perecíveis. No entanto, a rastreabilidade e monitoramento ao longo da cadeia tradicional frequentemente enfrentam desafios de falta de transparência e confiabilidade, levando a perdas e riscos à saúde pública. Este artigo propõe uma arquitetura que integra a Internet das Coisas (IoT), um banco de dados híbrido composto por Hyperledger e MongoDB, e computação em névoa e borda para aprimorar a logística desses produtos. A IoT possibilita a coleta de dados em tempo real sobre as condições ambientais, enquanto a computação em névoa e borda processa os dados mais próximos da fonte, permitindo ações instantâneas com atuadores. O banco de dados híbrido assegura a rastreabilidade ao armazenar informações de forma segura e imutável. Foi desenvolvido um protótipo que passou por testes, demonstrando os benefícios da abordagem híbrida de banco de dados e aprimorando a rastreabilidade e eficiência na gestão da cadeia de abastecimento. Embora desafios relacionados à escalabilidade tenham sido identificados, esta arquitetura possui o potencial de elevar os padrões do setor.

### Recursos principais:

A arquitetura proposta consiste em seis camadas: camada de usuário, camada de aplicação, camada data gateway, camada de banco de dados híbrido, camada de computação em névoa e camada de IoT com computação de borda.

- **Camada de usuário**: inclui produtores/fabricantes, empresas de logística, consumidores e administradores do sistema. Eles desempenham papéis específicos, como atribuir códigos eletrônicos aos produtos, monitorar parâmetros de armazenamento e gerenciamento do sistema.
- **Camada de aplicação**: oferece serviços como gerenciamento de usuários, dispositivos, produtos, regulamentação de qualidade e rastreabilidade de produtos.
A camada data gateway atua como middleware para processamento, armazenamento e consulta de dados. Ela inclui serviços como gerenciador de armazenamento, verificador de integridade de dados e análise de risco e alerta inteligente.
- **Camada de banco de dados híbrido**: combina blockchain e bancos de dados tradicionais para otimizar o desempenho e a escalabilidade. Os dados são armazenados on-chain e off-chain de acordo com a necessidade.
- **Camada de computação em névoa**: inclui gerenciamento de armazenamento, dispositivos, controle de conexão e validação de dados para dispositivos IoT.
- **Camada de IoT, juntamente com a computação de borda**: abrange o gerenciamento e monitoramento de dispositivos, o controle das conexões e a geração de alertas, bem como a capacidade de atuadores para realizar ações com base nas informações ambientais coletadas.

## Conclusão

Neste artigo, apresentamos uma arquitetura para a gestão de cadeias de suprimentos de produtos perecíveis. Essa arquitetura abrange várias camadas, desde a camada do usuário até a camada de IoT, e discutimos detalhadamente o papel de cada uma delas no sistema. 

Os testes realizados demonstraram que a abordagem híbrida de banco de dados, que combina MongoDB e Hyperledger Fabric, é eficaz. No entanto, também identificamos desafios relacionados ao dimensionamento na camada de borda e observamos tempos de inserção e consulta de dados que podem ser aprimorados. Com um planejamento adequado, esta solução pode melhorar a gestão de produtos perecíveis. 

Em trabalhos futuros, pretendemos explorar a integração de tecnologias avançadas, como Inteligência Artificial, para aprimorar a análise de dados e a validação em cenários do mundo real.
