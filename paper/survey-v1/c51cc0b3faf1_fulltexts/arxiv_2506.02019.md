## CHATCFD: AN LLM-DRIVEN AGENT FOR END-TO-END CFD AUTOMATION WITH STRUCTURED KNOWLEDGE AND REASONING

E Fan 2, † , Kang Hu 1, † , Zhuowen Wu 3 , Jiangyang Ge 1 , Jiawei Miao 3 , Yuzhi Zhang 3 , He Sun 4,5 , Weizong Wang 1,6 , and Tianhan Zhang 1,7,*

1 School of Astronautics, Beihang University

2 Department of Mechanics and Aerospace Engineering, Southern University of Science and Technology 3 DP Technology, Beijing

4 College of Future Technology, Peking University

5 National Biomedical Imaging Center, Peking University

6 State Key Laboratory of High-Efficiency Reusable Aerospace Transportation Technology

7 Key Laboratory of Spacecraft Design Optimization and Dynamic Simulation Technology, Ministry of Education † These authors contributed equally to this work and are listed alphabetically.

* Corresponding author, thzhang@buaa.edu.cn

## ABSTRACT

Computational Fluid Dynamics (CFD) is essential for advancing scientific and engineering fields, but is hindered by operational complexity, high expertise requirements, and limited accessibility. This paper introduces ChatCFD, a LLM-driven agent system for end-to-end CFD automation. Powered by DeepSeek-R1/V3, a multi-agent architecture, structured OpenFOAM knowledge bases, precise error locator, and iterative reflection, ChatCFD dramatically outperforms prior systems. On 315 benchmark cases it attains 82.1% execution success (vs. 6.2% MetaOpenFOAM, 42.3% FoamAgent) and, crucially, 68.12% physical fidelity -the first rigorous metric capturing whether a runnable simulation is scientifically meaningful. A dedicated Physics Interpreter achieves 97.4% summary fidelity , exposing the striking LLM gap between fluent narration and enforcing dozens of tightly coupled physical constraints in executable code. Resource-efficiency analysis highlights ChatCFD's practical advantage, achieving an average of 192.1k tokens and $0.208 per case -half the tokens of Foam-Agent and 1.5 × cheaper than MetaOpenFOAM. Ablation studies confirm structured domain knowledge and reasoning is indispensable: removing the Solver Template DB collapses accuracy to 48%, while the Error Locator Module proves the single most critical component. Flexibility experiments show autonomous solver selection across compressible/incompressible and steady/transient regimes (95.23% success) and turbulence-model switching (100% success), even on unseen configurations. By faithfully reproducing complex literature cases NACA0012 airfoil and supersonic nozzle with 60-80% end-to-end success where baselines fail entirely, ChatCFD establishes new benchmarks for AI-driven CFD automation. Physics coupling analyses reveal higher resource demands for multi-physics-coupled cases, while LLM bias toward simpler setups introduces persistent errors. ChatCFD's modular, MCP-compatible design directly enables collaborative multiagent networks and paves the way for scalable AI-driven CFD innovation. The code for ChatCFD is available at: https://github.com/ConMoo/ChatCFD .

## 1 Introduction

Computational Fluid Dynamics (CFD) is a cornerstone technology in diverse scientific and engineering disciplines, including aerospace [1, 2], energy systems [3, 4, 5], urban environment [6, 7, 8], combustion [9, 10, 11], astrophysics[12], and biomedical applications [13, 14]. It provides indispensable tools for simulating intricate fluid behaviors, thereby enabling design innovation and scientific discovery [15, 16]. However, CFD implementation faces substantial barriers, requiring deep domain expertise [17] and often relying on expensive commercial software. Even experts invest significant time in tasks such as solver selection, model setup, mesh generation, boundary condition definition, and post-processing [18]. These demands, combined with high computational costs, limit CFD accessibility for smaller organizations and hinder broader innovation [16]. Thus, there is an urgent need for automated, intuitive, and affordable CFD solutions. To meet this need, we present ChatCFD, an AI-driven agent system that streamlines CFD workflows, facilitating automated scientific discovery in fluid mechanics and engineering by enabling rapid iteration and exploration of complex flow phenomena without extensive manual intervention.

Recent progress in Artificial Intelligence (AI) has transformed the automation of sophisticated scientific processes. Large Language Models (LLMs), such as GPT [19], Gemini [20], and DeepSeek [21], alongside multi-agent frameworks like MetaGPT [22] and AutoGen [23], excel in natural language processing, code generation, and reasoning via chainof-thought techniques [24]. Frameworks like ReAct [25] synergize reasoning and acting in LLMs, while Reflexion [26] uses verbal reinforcement learning to improve agent performance through self-reflection. Toolformer [27] teaches LLMs to self-learn tool usage. Retrieval-Augmented Generation (RAG) [28] further bolsters these agents by incorporating domain-specific knowledge, reducing hallucinations, and improving domain adaptation for tasks like CFD automation.

This synergy has enabled LLM-based agents to automate CFD pipelines, from interpreting user inputs to running simulations. OpenFOAM [29], a leading open-source CFD platform, supports this trend due to its flexibility and community backing. Notable efforts include MetaOpenFOAM [30], which automates simulations from natural language using RAG and MetaGPT, though limited to tutorial-level cases and struggling with complex geometries; OpenFOAMGPT [31], evaluating RAG-augmented LLMs for zero-shot setup and boundary adjustments, but facing difficulties in turbulence model modifications; Foam-Agent [32], introducing hierarchical retrieval for dependencyaware generation, yet constrained in handling unseen configurations; and NL2FOAM [33], fine-tuning LLMs for natural language-to-CFD translation using LoRA [34], but reliant on domain-specific datasets and less effective for diverse physical models. Despite these advances, existing systems often fall short in end-to-end automation for complex cases without human oversight, particularly due to a fundamental challenge: LLMs lack sufficient training on niche scientific tasks like OpenFOAM setup, necessitating expert-designed architectures to bridge this gap through structured reasoning, reflection, and knowledge integration.

Existing CFD automation agents are often limited by their focus on rudimentary tasks, typically restricted to OpenFOAM tutorial-level examples, raising concerns about their generalization to complex, unseen cases. This challenge primarily stems from the scarcity of domain-specific scientific coding corpora, such as OpenFOAM setups, requiring sophisticated multi-agent system designs to effectively incorporate physical knowledge into agent architectures, thereby enhancing LLMs' capabilities in specialized scientific tool invocation. Recent representative work in this area includes SciToolAgent [35], which develops a knowledge graph-driven agent for integrating hundreds of scientific tools across biology, chemistry, and materials science, using graph-based RAG and safety-checking to automate workflows like protein engineering and chemical synthesis. While not directly focused on CFD, its approach to tool integration and domain adaptation offers valuable insights for fluid automation agents in handling coupled physics and complex simulations.

Three critical challenges in CFD automation agents remain unaddressed. First, incorporating domain-specific knowledge into CFD agents remains challenging. Engineering-scale simulations require comprehensive initial and boundary condition specifications, necessitating domain-specific structural thinking and seamless integration of specialized domain expertise into the automation pipeline. Second, designing effective agent frameworks, especially for future multi-MCP agent collaborations, is essential to maximize complementary capabilities. Typical issues include whether we need to contain mesh generation in the agent scope, where real-world geometries exceed basic OpenFOAM tools like blockMesh and require specialized integration of external meshes. Third, enabling large-scale testing and evolution of fluid agents is crucial. Current reliance on simplistic natural language descriptions often leads to ambiguous case setups and hinders scalability. In contrast, CFD literature, comprising millions of papers with detailed descriptions and standard results, remains underutilized, impeding agent evolution. These shortcomings highlight a critical gap in current CFD automation frameworks.

To overcome these challenges, we introduce ChatCFD, an automated agent system for OpenFOAM simulations, as shown in Figure 1. ChatCFD processes multi-modal inputs (e.g., research papers, meshes) through an interactive interface, utilizing LLMs ( DeepSeek-R1 , DeepSeek-V3 ), multi-agent architecture, and OpenFOAM knowledge to

Figure 1: Overview of the ChatCFD automated agent system for streamlining CFD simulations within the OpenFOAM framework. ChatCFD enables researchers and engineers to configure and execute simulations with minimal CFD or OpenFOAM expertise. The system comprises three core components: (1) an interactive chat interface for users to input case descriptions or upload mesh files, (2) a thinking system, the core decision-making module (detailed in Figure 2), and (3) a simulation engine that executes cases, collects error logs, and delivers final results.

<!-- image -->

manage the full workflow. It supports intricate setups via iterative refinement, handling diverse physical models and external meshes, surpassing prior agents in adaptability to unseen cases.

This paper details ChatCFD's architecture and validates its components experimentally. The structure is as follows: Section 2 describes the pipeline design, Section 3 presents results, and Section 4 summarizes findings. Appendices provide examples and prompts for operational clarity.

## 2 ChatCFD Pipeline Design

ChatCFD is an automated CFD agent system that leverages the OpenFOAM framework to process multi-modal user inputs, including research articles and mesh files, to configure and execute CFD simulations based on user instructions. As illustrated in Fig. 2, the ChatCFD framework implements a comprehensive workflow to streamline simulation setup, execution, and analysis.

## 2.1 Stage 0: Knowledge Base Construction

The preprocessing aims to establish a foundational knowledge base for CFD tasks, performed in advance to optimize ChatCFD's operations. Ideally, this knowledge base contains general CFD principles (e.g., numerical differentiation, discretization schemes), OpenFOAM-specific manuals, and a comprehensive example library. Fine-tuning large language models (LLMs) for general CFD knowledge is beyond the current scope; readers are referred to recent work by Dong et al. [33]. This study focuses on preprocessing OpenFOAM manuals and publicly accessible examples.

Through a systematic analysis of OpenFOAM manuals and tutorials, we construct a structured JSON database that records key parameters such as solver configurations, turbulence models, and file dependencies, laying the foundation for accurate case setup and error correction. Additionally, a vector database is established to further enhance the understanding of OpenFOAM and to generate more precise OpenFOAM case descriptions. As shown in Fig. 3, the

Figure 2: Architecture of the ChatCFD framework for automated CFD simulations, illustrating the four-stage workflow and agent structure. The stages are: (1) Knowledge Base Construction, creating a JSON database from OpenFOAM manuals and tutorials; (2) User Input Processing, enabling user interaction via natural language or document and mesh uploads; (3) Case File Generation, generating OpenFOAM case files using the knowledge base; and (4) Execution and Error Reflection, running simulations, converting meshes with fluentMeshToFoam , and resolving errors (dimension mismatches, missing files, persistent errors, general issues) using RAG-based modules ReferenceRetriever and ContextRetriever . The agent structure integrates DeepSeek-R1 and DeepSeek-V3 for intelligent processing, with iterative error correction.

<!-- image -->

database primarily comprises two components: a Structured Database in JSON format and a Vector Database for efficient retrieval. The Structured Database includes the following parts:

- File Dependency and Structure DB (db1) : Designed to construct a hierarchical file structure model. This model is based on OpenFOAM tutorial cases and categorized primarily by solver type and turbulence model to store the necessary file structure for specific combinations. It also records the actual file structure of each tutorial case, which is utilized to construct and validate the space of permissible configurations.
- Boundary Condition DB (db2) : Comprehensively catalogs all available boundary condition types within OpenFOAM. For each boundary condition, it precisely documents the required additional setting keywords and their configuration specifications, ensuring the completeness and validity of the boundary setup.
- Parameter Dimension DB(db3) : Systematically records the standard dimensional configuration rules for all physical quantities in OpenFOAM. This database further defines the dimensional requirements and consistency constraints for relevant field variables under varying physical problems or specific solving environments.
- Solver Template DB (db4) : Systematically stores the core file content and structured metadata of fundamental OpenFOAM tutorial cases. The metadata includes: the case's designated solver, turbulence model, and physical problem

Figure 3: Based on the knowledge base constructed from the contents such as the OpenFOAM manual and OpenFOAM tutorials, and how it functions within the ChatCFD framework

<!-- image -->

category. This information serves as the key basis for achieving high-accuracy classification and rapid screening of tutorial cases, thus providing precise reference templates for automated case initialization and configuration.

All configuration files are carefully parsed, converted to JSON format, and tagged with metadata (e.g., solver specifications and flow characteristics) for easy retrieval and application. OpenFOAM-specific header files are removed, and cases with external dependencies or auxiliary folders are excluded to maintain case autonomy and reduce the influence of irrelevant configuration elements on the LLM's response quality. Furthermore, for cases where non-uniform field definitions occupy a large volume of file content (potentially exceeding dialogue context limits), placeholders (e.g., \_\_CELL\_COUNT\_\_ , \_\_FACE\_COUNT\_\_ , \_\_SCALAR\_DATA\_PLACEHOLDER\_\_ ) are used to replace the non-uniform field data. This allows the LLM to focus on critical settings (such as boundary and initial conditions) while mitigating the issue of a reduced case count due to complex data.

Initially, preprocessed manual and tutorial data are used to define file dependency relationships, which are accessed during the RAG process to enhance the agent's error correction capabilities. To ensure robust operation, strategic case filtering is implemented by removing OpenFOAM-specific headers and excluding cases with external dependencies or auxiliary folders. This maintains case independence, reducing the risk of degraded LLM response quality due to irrelevant configuration elements. For cases where non-uniform field definitions occupy significant file content and may exceed the dialogue context limits of large models, placeholder tokens such as \_\_CELL\_COUNT\_\_ , \_\_FACE\_COUNT\_\_ , \_\_SCALAR\_DATA\_PLACEHOLDER\_\_ are used to replace non-uniform field data. This approach enables the LLM to focus on critical settings, such as boundary and initial conditions, while ignoring non-uniform field configurations. As a result, it mitigates the reduction in the number of usable OpenFOAM tutorial cases caused by complex non-uniform field data. Filtering criteria and implementation details are provided in Appendix A. Subsequently, the example library is preprocessed through systematic extraction and organization. All case configuration files are parsed and converted into a structured JSON format, tagged with metadata such as solver specifications, turbulence model types, and flow regime characteristics. Analysis of solver and turbulence model distributions determines specific file requirements for each configuration. For instance, the simpleFoam solver requires configuration files including system/fvSolution , system/controlDict , system/fvSchemes , 0/p , 0/U , constant/transportProperties , and constant/turbulenceProperties . Similarly, the kω SST turbulence model requires system/fvSolution , system/controlDict , system/fvSchemes , 0/p , 0/U , 0/k , 0/omega , 0/nut , constant/turbulenceProperties , and constant/transportProperties .

## 2.2 Stage 1: User Input Processing

The initial stage of ChatCFD features an interactive, multi-modal interface that enables users to define CFD simulations by either conversing with the DeepSeek-R1 model or uploading documents and mesh files, built using the Streamlit Python framework. This interface leverages the preprocessed knowledge base to facilitate accurate case specification

with minimal CFD expertise. The knowledge base, comprising OpenFOAM manuals, tutorial cases, and a structured JSON database of solver configurations and turbulence models, informs the system's natural language processing and case extraction, ensuring robust and context-aware interactions. The workflow comprises four key phases:

- Input Submission : Users can upload research articles or describe the case through dialogue to provide the basis for case analysis.
- Case Extraction : The system extracts and catalogs CFD cases from uploaded documents or conversation inputs, identifying solver configurations and physical models, and presents them in a structured format with unique identifiers (e.g., Case 1, Case 2).
- Case Selection : Through natural language interaction, the system guides users to select a target case, provides detailed specifications for verification, and confirms the selection.
- Mesh Integration : The system assists users in uploading mesh files, currently supporting Fluent format .msh files. For OpenFOAM default meshes, users must specify the path to the corresponding polyMesh directory.

Upon completing this workflow, ChatCFD aggregates three essential components for downstream processing, as shown in Fig. 2:

- Case Documentation : The system retains the full research article or conversation details, providing access to all technical specifications, including solvers, boundary conditions, and simulation parameters.
- Case Specification : The selected case is recorded with detailed metadata, including solver types, numerical schemes, physical models, and source references, ensuring precise setup.
- Mesh Data : Robust file transfer and path management ensure reliable delivery of mesh data to subsequent stages.

## 2.3 Stage 2: Case File Generation

The CFD engineer layer, powered by the DeepSeek-R1 and DeepSeek-V3 models, initializes OpenFOAM case files through a streamlined three-phase process that leverages the preprocessed knowledge base.

## 2.3.1 Phase 1: File Configuration Synthesis

The system analyzes the case description from Stage 1 (User Input Processing) to identify suitable solvers and models, thereby generating a list of required case file configurations. This involves explicitly matching dependencies for solvers, models, and related settings to facilitate accurate file configuration. The workflow proceeds as follows:

- Relevant cases are retrieved based on the solver and model, yielding a corresponding file configuration list from which non-essential files (e.g., those for visualization or mesh generation) are removed. These refined cases are then input to the DeepSeek-V3 model to produce an optimized file configuration list.
- In the absence of direct matches, cases are aligned by solver type. For additional models, such as turbulence models, necessary files are inferred from analogous cases, synthesized into a new configuration list, verified for solver compatibility, and finalized via DeepSeek-V3 .

This approach integrates solver and physical model structures with case synthesis, enhancing initialization precision, minimizing file omissions, and substantially improving overall case file quality.initialization.

## 2.3.2 Phase 2: Case Setup Extraction

Detailed case configurations are derived from three primary sources: the file configuration list, user prompts, and the case library. The configuration list prevents omissions or redundancies in the case structure, while user prompts enable filtering of segmented case description chunks using the sentence-transformer/all-mpnet-base-v2 library and a similarity threshold. The DeepSeek-R1 model implements a hierarchical three-step extraction:

- Mesh Boundary Condition Extraction: Boundary names and types are extracted via regular expressions and tools like pyFoam, forming a persistent dictionary for subsequent physical field generation and error handling.
- Boundary Condition Validation: Boundary conditions are identified from filtered chunks, checked for compatibility with OpenFOAM-v2406, and subjected to rigorous naming and format validation.
- Physical Field File Value Setup: Field files (e.g., 0/p , 0/nut , 0/nutTilda ) are identified from the configuration list. Using prior boundary data, OpenFOAM-compliant templates are created, encompassing internal and boundary field settings differentiated by boundary names and associated keywords.

## 2.3.3 Phase 3: Configuration Validation and Refinement

For unspecified parameters, such as discretization schemes and solver algorithms in system/fvSolution and system/fvSchemes , the DeepSeek-R1 model applies advanced reasoning grounded in CFD best practices and case-specific physics to formulate coherent configurations. This maintains parameter interdependencies and ensures alignment with the outlined flow dynamics and solver demands.

A comprehensive validation and correction mechanism further refines the generated files. The system scrutinizes structures in the system and constant directories, confirming proper dependencies and completeness through crossreferences with library exemplars. Particular emphasis is placed on physical field dimensions to mitigate LLM limitations in interpreting such constraints. For example, variables like p and alphat -which differ dimensionally between compressible and incompressible flows-are calibrated according to case type; ambiguous cases retain initial values. This process curtails errors like dimensional inconsistencies, bolsters initial file robustness, and diminishes iterative corrections in Stage 3.

## 2.4 Stage 3: Error Correction and Reflection

Stage 3 of the ChatCFD pipeline employs a modular architecture for automated error correction and iterative refinement in OpenFOAM case configurations. The following modules facilitate robust error handling:

- ReferenceRetriever : This module retrieves reference files from the preprocessed OpenFOAM knowledge base, matching solver and model specifications from benchmark tutorial cases. If no exact match is found, it prioritizes solver-compatible files. To balance guidance and prompt complexity, two reference files are selected to aid error detection and correction, leveraging the LLM's few-shot learning to enhance ChatCFD's accuracy in error locator and resolution.
- ContextRetriever : This module compiles current case configurations into a structured JSON format, providing detailed file content and modification trajectories over a specified period, tailored to other modules' needs. It supports targeted reflection and correction, particularly for errors involving physical coupling (e.g., pressure-density interactions).
- Error Locator Module : Comprising two components: (1) DeepSeek-V3 rapidly identifies suspicious files based on error messages, reducing the analysis burden; (2) DeepSeek-R1 performs detailed reasoning to pinpoint erroneous files by comparing their content against the case library, leveraging its advanced inference capabilities for precise error location.
- Error Correction Module : This module addresses file interdependencies in OpenFOAM through a three-step process: (1) DeepSeek-V3 identifies files requiring coordinated modifications to resolve coupling-related errors; (2) DeepSeek-R1 proposes actionable corrections by integrating error details, related file content, simulation requirements, and benchmark tutorial cases; (3) DeepSeek-V3 applies these corrections, ensuring compliance with OpenFOAM's file format standards.
- Reflection Module : Activated during persistent errors, this module collects error messages and file modification histories, formulating reflections in the format: 'For situation A, I considered B but overlooked C; next, I will apply D to resolve the issue.' These reflections are stored as reflection blocks within a &lt;reflexion&gt; tag in the reflection history, enhancing iterative error correction by integrating with error locator and correction modules.

Errors detected during simulation are categorized into two types:

- General Errors : These, comprising over 70% of runtime errors, include incorrect keywords, formatting issues, or floating-point errors. Correction involves a streamlined workflow: suspicious file detection ( DeepSeek-V3 ), error confirmation ( DeepSeek-R1 ), retrieval of related tutorial files ( ReferenceRetriever ), proposal of corrections, and file modification. For missing file errors, DeepSeek-V3 directly generates the required file using reference cases, followed by dimensional validation to ensure quality.
- Persistent Errors : ChatCFD leverages short-term memory (recent error messages and file modification histories) and long-term memory (reflection histories) to address recurring issues. Reflection histories, stored as structured insights, enhance the system's ability to adapt and resolve complex errors iteratively.

## 2.5 Stage 4: Post-Processing

To enhance the transparency and user intelligibility of the computation process, ChatCFD integrates a dedicated Physics Interpreter in the fourth stage of the workflow. This component, implemented using the DeepSeek-R1 model, automatically produces a technical summary immediately after the case is successfully generated and deemed

executable. The output is presented in Markdown format, supplemented by structured expressions such as tables and file tree diagrams. Its content systematically covers multiple dimensions, including fundamental information, physical problem background, special configuration settings, file structure, and key file contents, thereby comprehensively ensuring the transparency of the simulated physical process for the user.

## 2.6 LLMArchitecture and Configuration

ChatCFD employs a dual-model architecture, integrating the DeepSeek-R1 and DeepSeek-V3 large language models (LLMs) to execute a comprehensive CFD workflow for case generation and error correction. Each model is strategically deployed based on its strengths to optimize performance across Stages 2 and 3.

The DeepSeek-R1 model, with its superior text comprehension and reasoning capabilities, is utilized for complex, high-stakes tasks. In Stage 2, it extracts critical simulation parameters (e.g., boundary conditions, initial conditions) from case descriptions and generates complete initial case files. In Stage 3, it pinpoints error-causing files, proposes detailed correction strategies, and conducts reflection on persistent errors. However, DeepSeek-R1 's propensity for hallucinations-such as inserting unnecessary Markdown formatting, which can disrupt error correction-poses challenges. Despite mitigation efforts using prompt engineering and pydantic for structured outputs, these issues persist. Additionally, its inference time increases significantly with longer prompts, making it less efficient for routine tasks.

The DeepSeek-V3 model excels in instruction following and rapid response, making it ideal for structured, lowcomplexity tasks. In Stage 2, it performs quick file structure validation, while in Stage 3, it identifies suspicious or missing files and generates corrected file content based on provided instructions. To address DeepSeek-V3 's limited reasoning capacity, tasks are simplified into structured formats (e.g., providing configuration file lists to narrow analysis scope), and the model is prompted to justify its actions (e.g., explaining why a file is suspicious). This enhances its effectiveness while maintaining computational efficiency. This dual-model approach leverages DeepSeek-R1 's reasoning for complex CFD tasks and DeepSeek-V3 's efficiency for structured operations, ensuring robust performance across diverse case types (e.g., benchmark tutorial cases, literature-derived cases) while minimizing computational costs.

ChatCFD maintains a comprehensive log for each CFD case, capturing all question-answer interactions and actions taken during reflection-based error correction. These logs are invaluable for analyzing the behavior of the DeepSeek-R1 and DeepSeek-V3 models within the OpenFOAM framework, particularly in identifying limitations in handling complex CFD configurations. The recorded reflection iterations, which document corrective actions for persistent errors, provide critical insights into scenarios where ChatCFD's performance is suboptimal, guiding targeted improvements in the pipeline. Additionally, a concise summary of reflections for recurring errors is maintained, offering a deeper understanding of ChatCFD's behavior in complex scenarios. These logs have the potential to form a novel database, serving as a valuable resource for enhancing LLM-based agents in addressing high-complexity CFD problems within the MCP framework.

The experimental validations described in Section 3 were conducted on the VolcEngine platform [36], which provides the pricing structure for token usage and computational costs. For DeepSeek-R1 , input tokens cost $0.00055 per 1,000 (0.004 RMB), and output tokens cost $0.0022 per 1,000 (0.016 RMB). For DeepSeek-V3 , input tokens cost $0.00021 per 1,000 (0.0015 RMB), and output tokens cost $0.00082 per 1,000 (0.006 RMB). All token consumption and cost metrics reported in Section 3 are based on this pricing, enabling precise evaluation of ChatCFD's computational efficiency across benchmark tutorial cases, perturbed variant cases, and literature-derived cases.

## 3 Results and discussion

The performance of ChatCFD was rigorously evaluated through a series of validation experiments encompassing three distinct categories of CFD cases: (i) 205 benchmark tutorial cases drawn from OpenFOAM tutorials and the OpenFOAM wiki, serving as standardized references for foundational validation; (ii) 110 perturbed variants, derived by systematically altering key parameters (e.g., boundary conditions, solver settings, or physical properties) in the benchmark cases to assess robustness and sensitivity; and (iii) 2 advanced literature-derived cases, directly prompted from published research articles to test real-world applicability in complex, unseen scenarios. ChatCFD demonstrated an operational success rate of 82.1% across the 315 benchmark and perturbed cases, with success defined as error-free configuration and execution leading to converged simulations. For the advanced literature-derived cases, success rates ranged from 60% to 80%, reflecting the challenges of handling intricate, domain-specific configurations without prior exposure. These results represent a substantial advancement in LLM-driven CFD automation, particularly in reducing the expertise barrier for setup and execution within the OpenFOAM framework.

## 3.1 Benchmarking Execution Success Rates Across Agents

To evaluate ChatCFD, a comprehensive set of test cases was curated from OpenFOAM tutorials and the OpenFOAM wiki. Cases involving complex mesh generation techniques were filtered out, resulting in 205 benchmark tutorial cases. From these, 11 cases were selected and systematically perturbed by modifying their descriptions and adjusting parameters (e.g., boundary conditions, solver settings, or physical properties), yielding an additional 110 perturbed variant cases. This approach ensures both foundational validation and robustness testing under varied conditions. Figure 4 illustrates the performance of ChatCFD compared to MetaOpenFOAM and Foam-Agent across these cases, showing success rates by category, case distribution, and overall performance. Success is defined as the generation of error-free case files leading to converged simulations. As shown in Figure 4(c), ChatCFD achieves an overall success rate of 82.1%, significantly outperforming MetaOpenFOAM (6.2%) and Foam-Agent (42.3%). Figure 4(a) highlights ChatCFD's consistent advantage across case categories, demonstrating its superior performance and generalization capability.

Figure 4: Comparison of success rates across three CFD agents (ChatCFD, MetaOpenFOAM, Foam-Agent) for 205 benchmark and 110 perturbed OpenFOAM tutorial cases. (a) Success rates by case category. (b) Distribution of test cases across categories. (c) Overall success rate comparison.

<!-- image -->

Notably, all agents exhibit lower success rates in combustion and multiphase flow categories due to their interdisciplinary complexity, which involves intricate file structures and content. For example, in the combustionFlame2D case, the DeepSeek-R1 model always erroneously formulates the thermo.compressibleGas file by adding a mixture dictionary for species like O2 and H2O , leading to persistent structural errors. This stems from an intuition that the gas is a mixture, overlooking OpenFOAM's specific requirement that thermo.compressibleGas defines compressible gas thermophysical properties, while thermophysicalProperties is the appropriate file for mixture settings. To address this, ChatCFD employs two strategies: (1) retrieving and integrating relevant OpenFOAM case libraries, leveraging the LLM's few-shot learning to adapt to specific file structures; and (2) enhancing the reflection module, which enables the system to review decision trajectories after repeated errors, identify knowledge or action discrepancies, and adopt corrected solutions. These mechanisms significantly improve ChatCFD's handling of challenging cases, boosting success rates in complex scenarios.

The performance of ChatCFD, MetaOpenFOAM, and Foam-Agent was evaluated through detailed metrics on 205 benchmark tutorial cases from OpenFOAM, focusing on token consumption, reflection iterations, and computational

Figure 5: Performance statistics for different agents across 205 benchmark tutorial cases. (a) Average token consumption per case. (b) Distribution of token consumption. (c) Average number of reflection iterations per case. (d) Average token consumption per reflection iteration. (e) Average computational cost (in monetary terms). (f) Distribution of reflection iteration ratios, excluding zero and limit-reaching cases.

<!-- image -->

cost. Figure 5 presents a comparative analysis of these metrics, highlighting ChatCFD's efficiency and robustness in automated CFD workflows. As shown in Figure 5(a), ChatCFD exhibits the lowest average token consumption per case, approximately half that of Foam-Agent, with MetaOpenFOAM ranking second. This efficiency underscores ChatCFD's high success rate combined with reduced costs, making it practical for large-scale CFD applications. While ChatCFD and MetaOpenFOAM have comparable total token consumption, Figure 5(e) reveals a significant cost disparity, with ChatCFD's computational cost being approximately 1.5 times lower than MetaOpenFOAM's. This is primarily due to ChatCFD's strategic use of the cost-efficient DeepSeek-V3 model (with a cost half that of DeepSeek-R1 ) for specific pipeline stages, as illustrated in Figure 5(b). In contrast, MetaOpenFOAM's lower success rate necessitates frequent error corrections, generating additional file content and inflating costs. Figure 5(c) shows that MetaOpenFOAM requires the highest average number of reflection iterations per case, followed by Foam-Agent, with ChatCFD requiring the fewest. However, Figure 5(d) indicates that MetaOpenFOAM consumes fewer tokens per reflection iteration due to its simpler error-correction approach, which dilutes the high token cost of initial case setup across multiple iterations. Conversely, ChatCFD and Foam-Agent employ more sophisticated error-reflection strategies. Foam-Agent, for instance, maintains an 'error correction trajectory' that logs modifications to file i at iteration j , supplemented by contextual data from related files and reference cases. This increases its per-iteration token consumption. ChatCFD, while similarly providing contextual data, optimizes efficiency by maintaining both short-term modification trajectories and long-term reflection summaries. This dual approach enables precise filtering of irrelevant information and concise summarization of critical context, significantly reducing token consumption during LLM interactions compared to Foam-Agent, as evidenced in Figure 5(d). Figure 5(f) depicts the distribution of reflection iteration ratios, defined as the number of reflections divided by the maximum allowed iterations, reflecting each agent's assessment of case complexity. The x-axis represents the reflection ratio, and the y-axis indicates the proportion of cases at each complexity level. Higher curves indicate stronger error-correction capabilities at a given reflection ratio. All three agents exhibit a rapid decline in solved cases followed by saturation, suggesting that increasing reflection iterations beyond a certain point yields diminishing returns. For Foam-Agent, with a reflection limit of 49 iterations, the curve flattens at approximately 50% (around 25 iterations), indicating that additional reflections beyond this threshold rarely resolve errors. ChatCFD's curve consistently lies above the others, demonstrating more effective reflections that are likelier to correct errors. However, even for ChatCFD, reflections beyond 75% of the limit (approximately 22 iterations) often fail to resolve remaining errors, highlighting a current limitation in LLM-based CFD agents where complex case errors require advanced knowledge integration or alternative strategies beyond iterative reflection.

## 3.2 Physical Fidelity Evaluation and Semantic Error Analysis

In Section 3.1, we assessed ChatCFD, MetaOpenFOAM, and Foam-Agent using success rate as the primary metric. While this allows rapid quantitative evaluation, successful execution does not ensure physical meaningfulness in CFD. To address this, we introduce physical fidelity as a core metric, evaluating whether executed simulations accurately capture user-specified physics. We also re-examine the link between run success and physical fidelity via agent-generated cases.

We define physical fidelity through a three-tier protocol. A case qualifies only if it meets all criteria without simplifications that compromise the scientific objective or user requirements:

- Conditions Verification : Initial and boundary conditions match the user-described scenario.
- Model Confirmation : Physical properties, models, and parameters are appropriate and consistent.
- Flow Feature Comparison : Post-processed visualizations align with ground-truth flow characteristics.

Using this protocol, we analyzed ChatCFD and Foam-Agent results (Fig. 6). ChatCFD's physical fidelity rate among runnable cases (phy|run) is lower than Foam-Agent's, but its overall rate is higher, confirming superior performance. Note that phy|run, limited to runnable cases, may inflate for agents handling only simple setups.

Figure 6: Physical fidelity evaluation of ChatCFD. (a) Success rates of ChatCFD and Foam-Agent under three criteria: run - execution success (simulation completes without crashing); phy | run - physical fidelity conditioned on execution success (proportion of runnable cases that are physically correct); phy - overall physical fidelity (proportion of all cases that are both runnable and physically meaningful). (b) Root-cause analysis of semantic failures in ChatCFD-generated cases that execute successfully but lack physical fidelity, categorized by error type (boundary/initial conditions, physical properties/models, numerical schemes, etc.).

<!-- image -->

Additionally, we summarized the causes of physical incorrectness in runnable cases. The results, shown in Fig 6, indicate that incorrect settings of boundary conditions, initial conditions, and initial fields in multiphase flows are the main contributors, accounting for 60 % of the issues, as these require deep physical insight and geometric awareness. Such misconfigurations rarely trigger convergence failures, allowing apparent success that misleads agents. Similarly, omissions in setFieldsDict often evade detection, exacerbating infidelity.

To further enhance transparency, we integrated a Physics Interpreter, evaluated using DeepSeek-V3 and GPT-4o for summary rationality. It achieved 97.4% fidelity, confirming accurate natural-language conveyance of the agent's intended physics.

The contrast between the near-perfect fidelity of the Physical Interpreter (97.4%) and the overall Physical Fidelity of only 68.12% is particularly interesting. This suggests that while a seemingly simple user request (e.g., "2-D flow around a cylinder at Re = 100") can be almost perfectly restated in natural language by the LLM, its faithful implementation in executable OpenFOAM code requires dozens of tightly interdependent, domain-specific settings (such as inlet velocity profile, turbulence suppression, kinematic viscosity, blockage ratio, spanwise boundary conditions, etc.). This highlights a key limitation of current LLMs: they excel at high-level linguistic narration but consistently struggle to enforce the complete chain of physical constraints within the generated executable code.

## 3.3 Ablation Study on ChatCFD Architectural Components

To quantitatively assess the contribution of the proposed database system and functional modules to the overall system performance, we conducted an ablation study using 20% of the cases sampled from the OpenFOAM benchmark test set. It should be noted that the degree of influence exerted by different components will vary due to the inherent differences in task types. Consequently, for the specific task of literature replication, we will perform additional ablation experiments in Section 3.5.1 for further, dedicated analysis.

The first four experiments were designed to evaluate the impact of the knowledge databases on ChatCFD's performance. This knowledge base primarily influences Stage2 and Case File Generation and Error Correction and Reflection in Stage3. Simultaneously, the experimental groups "w/o mod1" (removal of the Error Locator Module) and "w/o mod2" (removal of the Reflection Module) were used to verify and quantify the efficacy of the error locator and reflection mechanisms, with their influence predominantly focused on Error Correction and Reflection in Stage3.

The impact of the first four knowledge base ablation experiments on ChatCFD's performance, ranked from highest to lowest influence, is as follows: Solver Template DB (db4, 48 %), Boundary Condition DB (db2, 57 %), Parameter Dimension DB (db3, 68 %), and File Dependency and Structure DB (db1, 73 %). The results show that the Solver Template DB is utilized across all operational stages of ChatCFD and is crucial for leveraging the Large Language Model's few-shot learning capability to generate initially correct case files.

Figure 7: Ablation study quantifying the contribution of individual prior-knowledge components and core modules to ChatCFD's configuration accuracy. w/o File Dependency DB (DB1) : removes the structured JSON database that enforces mandatory file hierarchy and dictionary dependencies. w/o Boundary Condition DB (DB2) : disables the boundary-condition template database (no validated keyword requirements). w/o Parameter Dimension DB (DB3) : disables automatic dimension enforcement and file-format checks. w/o Solver Template DB (DB4) : removes all reference tutorial snippets and metadata. w/o Error Locator Module (Module 1) : disables the error-positioning module. w/o Reflection (Module 2) : disables the iterative reflection loop. w/ all : full ChatCFD system (all four prior-knowledge databases + all function modules).

<!-- image -->

Fig 7 shows that the removal of the Error Locator Module (w/o mod1) resulted in the lowest success rate. Here, we bypassed it and prompted DeepSeek-R1 to modify files based solely on error messages. The experimental results reveal two findings: 1) Accurate identifying of the erroneous file is a critical step for enhancing the agent's efficiency during the reflection and correction phase (Stage 3); and 2) Relying only on error messages is ineffective. Its shortcomings are particularly pronounced when dealing with errors that do not explicitly point to the source file, such as dimensional errors. Removing the Reflection Module (w/o mod2) causes the smallest performance drop because it primarily benefits highdifficulty cases, which constitute only a minor fraction of the benchmark. In contrast, the other components-knowledge bases and Error Locator-impact cases across all difficulty levels.

## 3.4 Flexibility Evaluation of Alternative OpenFOAM Configurations

To evaluate ChatCFD's ability to adapt to different physical assumptions and numerical settings, we designed two targeted experiments:

- Solver selection across regimes -Given only a textual description, ChatCFD autonomously determines the flow regime (compressible/incompressible, steady/transient) and selects the appropriate solver.
- Turbulence-model switching -With mesh and boundary conditions fixed, only the turbulence closure is changed, requiring automatic rewriting of constant/turbulenceProperties and related dictionaries.

Fig 8 summarizes the results of Experiment 1. ChatCFD correctly identified the required solver in 20 of 21 cases, including the unseen 'simpleCar' compressible-transient configuration (not present in the database). The sole failure occurred on 'buoyantCavity' under compressible-transient assumptions, which currently exceeds ChatCFD's supported scope. In the turbulence-switching experiment (Figure 9), ChatCFD successfully applied kϵ , kω SST, and laminar closures to every case-including the originally laminar damBreak -demonstrating robust handling of non-default configurations. These results confirm ChatCFD's strong generalization: it is not limited to tutorial-default settings but can reason about physically valid alternatives, significantly enhancing its practical utility for practical CFD exploration.

Figure 8: Solver flexibility under varying physical regimes. Blue circles indicate the reference (tutorial) setting; green checks mark successful execution with the identified solver; red crosses denote failure. ChatCFD correctly selects the appropriate solver in 20/21 of cases, including unseen compressible-transient scenarios.

| Case Name     | Incompressible &Steady-state   | Compressible &Steady-state   | Incompressible &Transient   | Compressible &Transient   |
|---------------|--------------------------------|------------------------------|-----------------------------|---------------------------|
| cavity        |                                |                              |                             |                           |
| pitzDaily     |                                |                              |                             |                           |
| airFoil2D     |                                |                              |                             |                           |
| cylinder      |                                |                              |                             |                           |
| simpleCar     |                                |                              |                             |                           |
| hotRoom       |                                |                              |                             |                           |
| buoyantCavity |                                |                              |                             | ×                         |

Figure 9: Turbulence-model flexibility. ChatCFD successfully switches between kϵ , kω SST, and other closures for all tested cases, including the originally laminar damBreak (100% success). Symbols follow Figure 8.

<!-- image -->

| CaseName                      | Turbulence model   | Turbulence model             | Turbulence model             | Turbulence model             | Turbulence model   |
|-------------------------------|--------------------|------------------------------|------------------------------|------------------------------|--------------------|
| cavity                        | kEpsilon           | RNGkEpsilon                  | kOmegaSST                    |                              |                    |
| pitzDaily                     | kOmegaSST          | SpalartAllmaras dynamicKEqnk | SpalartAllmaras dynamicKEqnk | SpalartAllmaras dynamicKEqnk |                    |
| squareBendLiq                 | kEpsilon           | kOmegaSST                    | SpalartAllmaras              |                              |                    |
| backwardFacingStep2D kEpsilon |                    | RNGkEpsilon                  | kOmegaSST                    |                              |                    |
| aeroilNACA0012                | RNGkEpsilon        | kOmegaSST                    | SpalartAllmaras              |                              |                    |
| Tjunction                     | kEpsilon           | kOmegaSST                    | realizableKE                 |                              |                    |
| dambreak                      | kEpsilon           | RNGkEpsilon                  | kOmegaSST                    |                              |                    |

## 3.5 Evaluation of CFD Cases Reproduced from Literature

To evaluate ChatCFD's ability to handle complex, real-world scenarios, two representative literature-derived cases were selected: an incompressible flow case and a compressible flow case, as illustrated in Figure 10. These cases were chosen for their comprehensive documentation in the source literature, providing detailed specifications for solver

<!-- image -->

(t) Boundary conditions type of the Naca0012 case provided by the reference.

(d) Boundary conditions type of the Nozzle case provided by the reference.

|                 | Airfoil                                                     | Inlet/outlet                                                 | Front and back          |
|-----------------|-------------------------------------------------------------|--------------------------------------------------------------|-------------------------|
| U p nuTilda nut | fixed-value  zero gradient fixed-value nutLowReWallFunction | freestream velocity freestreamPressure freestream freestream | empty empty empty empty |

Figure 10: Visualization of two literature-derived CFD cases. (a) Computational mesh for the NACA0012 airfoil case [37]. (b) Computational mesh for the Nozzle case [38]. (c) Boundary condition types for the NACA0012 case. (d) Boundary condition types for the Nozzle case.

|           | Pressure (P)     | Velocity (U)     | Temperature (T)   |
|-----------|------------------|------------------|-------------------|
| Inlet     | totalPressure    | zeroGradient     | totalTemperature  |
| Wall      | zeroGradient     | noSlip           | zeroGradient      |
| Far Field | waveTransmissive | waveTransmissive | waveTransmissive  |
| Outlet    | waveTransmissive | waveTransmissive | waveTransmissive  |

configurations, turbulence models, and mesh characteristics critical for accurate CFD reproduction. The incompressible flow case, based on Sun et al. [37], involves a NACA0012 airfoil at a 10° angle of attack, simulated using the simpleFoam solver with the Spalart-Allmaras turbulence model. The compressible flow case, drawn from Yu et al. [38], models a nozzle with a pressure ratio of 3, employing the rhoCentralFoam solver and the Spalart-Allmaras model.

Given reported challenges in modifying OpenFOAM turbulence models using large language models (LLMs), as noted by Pandey et al. [31], additional experiments were conducted to explore robustness across various turbulence models (e.g., kϵ , kω SST). Table 1 provides a comprehensive summary of all case configurations and experimental iterations. High-fidelity computational meshes were generated to ensure accuracy. The experimental protocol was structured as follows: Cases 1 (NACA0012) and 5 (nozzle) underwent an extensive model ablation study, with 10 iterations per case across five system configurations to assess sensitivity to model variations. Cases 2 through 4, using different turbulence models, were evaluated based on the full ChatCFD system configuration, with 10 iterations each, to focus on performance under optimal settings. This rigorous methodology resulted in 130 experimental runs, forming a robust dataset for analyzing ChatCFD's performance in complex, literature-derived scenarios. The results, detailed in subsequent analyses, demonstrate ChatCFD's capability to adapt to sophisticated CFD configurations, complementing its strong performance on benchmark tutorial and perturbed variant cases.

## 3.5.1 Ablation Study on ChatCFD Configurations for Literature Reproduction

To systematically dissect the contributions of ChatCFD's core components, ablation experiments were conducted across five distinct system configurations. These configurations vary in the complexity of case file generation (Stage 2, Section 2.3), the error-handling logic, and the integration of retrieval modules ( ReferenceRetriever and ContextRetriever ) in Stage 3 (Correcting Case Files). The configurations are defined as follows:

Table 1: Summary of Test Cases for ChatCFD Validation

| Case ID   | Case Mesh   | Physical Models                                      | Experimental Design                         |
|-----------|-------------|------------------------------------------------------|---------------------------------------------|
| Case 1    | NACA0012    | simpleFoam with Spalart-Allmaras                     | 10 runs × 5 system configurations           |
| Case 2    | NACA0012    | simpleFoam with k- ω SST                             | 10 runs × the complete system configuration |
| Case 3    | NACA0012    | simpleFoam with k- ϵ                                 | 10 runs × the complete system configuration |
| Case 4    | NACA0012    | simpleFoam with RNGk- ϵ                              | 10 runs × the complete system configuration |
| Case 5    | Nozzle      | rhoCentralFoam with hePsiThermo and Spalart-Allmaras | 10 runs × 5 system configurations           |

- Configuration A (Baseline System) : Employs a simplified parameter extraction process in Stage 2, bypassing the three-step hierarchical extraction. In Stage 3, error correction is limited to a basic 'general error' pathway without ReferenceRetriever (tutorial-based retrieval).
- Configuration B (Baseline with ReferenceRetriever ) and Simplified Setup Extraction : Retains the simplified Stage 2 extraction but enhances Stage 3 with the ReferenceRetriever module, enabling tutorial-based error correction within the general error pathway.
- Configuration C (Enhanced Setup and Reference Retrieval ) : Implements the full three-step setup extraction and configuration validation in Stage 2. Stage 3 includes general error correction augmented by the ReferenceRetriever module for improved accuracy.
- Configuration D (Full Error Reflection with Simplified Extraction) : Uses the simplified Stage 2 extraction but incorporates both general and persistent error correction modules in Stage 3, supported by ReferenceRetriever and ContextRetriever for comprehensive error handling.
- Configuration E (Complete System) : Represents the full ChatCFD pipeline, featuring comprehensive case setup extraction and validation in Stage 2 and complete error reflection modules in Stage 3, integrating both ReferenceRetriever and ContextRetriever for optimal performance.

Figure 11 illustrates ChatCFD's performance across Configurations A through E for two literature-derived cases: NACA0012 (Case 1) and Nozzle (Case 5). Two metrics were evaluated: (1) the 10-step operational success rate , defined as the ability to complete 10 simulation steps (time steps for transient cases or iterations for steady-state cases) without critical errors, and (2) the accurate configuration success rate , indicating exact compliance with the solver, turbulence model, and boundary condition specifications in the source articles. The 10-step criterion enables rapid assessment of operational stability, particularly for the reflection mechanism, as literature often lacks specified simulation durations. Configurations achieving this milestone were generally free of major operational faults, though minor deviations in boundary condition parameters could persist. Cases meeting the accurate configuration criterion exhibited simulation results in strong agreement with published data, as detailed in the appendix. These results underscore ChatCFD's robustness in handling complex, literature-derived configurations, complementing its performance on benchmark tutorial and perturbed variant cases.

For the NACA0012 Case 1, the data in Figure 11 show a clear progression in the 10-step operational success rate with increasing system sophistication. The baseline Configuration A achieved a 20% success rate. This trend highlights the positive impact of incrementally integrating retrieval modules and OpenFOAM-specific knowledge. However, despite its high operational success, Configuration D failed to achieve any accurate configuration success, unable to precisely replicate the case setup from the source article. This crucial gap was addressed by Configuration E, which incorporates comprehensive article interpretation and setup extraction as part of its advanced Stage 2 processing. With Configuration E, the accurate configuration success rate for the NACA0012 case improved significantly to 40%, while maintaining the 80% operational success rate. Analysis of the DeepSeek-R1 model's internal reasoning revealed that, without robust article interpretation and strict compliance with documented specifications in Stage 2, the LLM often defaults to generic boundary conditions and simplified parameters, such as standard inlet -type conditions for inflow velocities or pressureOutlet -type conditions for outflow pressures, instead of the specialized freestream conditions required for the NACA0012 airfoil. These erroneous configurations underscore the critical role of comprehensive case setup

Figure 11: Success rates of ChatCFD across five system configurations (A-E) for literature-derived cases: NACA0012 (Case 1) and Nozzle (Case 5). The horizontal axis represents configurations, and the vertical axis quantifies success rates. The 10-step operational success rate indicates completion of 10 simulation steps (time steps for transient cases or iterations for steady-state cases) without critical errors. The accurate configuration success rate reflects precise adherence to source article specifications.

<!-- image -->

extraction and validation in Stage 2, which Configuration E leverages to ensure fidelity to complex CFD requirements, enhancing ChatCFD's performance in complicated scenarios.

For the Nozzle Case 5, Figure 11 demonstrates a significant enhancement in the 10-step operational success rate with the integration of the ContextRetriever module. Configurations A through C, lacking ContextRetriever , achieved a 0% success rate for this metric, defined as completing 10 simulation steps (time steps for transient cases or iterations for steady-state cases) without critical errors. Configuration D, incorporating ContextRetriever , improved this rate to 60%. This advancement stems from ContextRetriever 's ability to diagnose complex coupling-related errors within the MCP framework, addressing a key limitation of earlier configurations. For example, errors in thermophysical parameters, such as compressibility ( Ψ ), were often misattributed to the constant/thermophysicalProperties file in prior configurations. In reality, these errors frequently arise from inconsistencies in equation of state specifications within boundary condition entries in field files (e.g., 0/p , 0/T ). By analyzing all relevant directories ( 0/ , constant/ , system/ ), ContextRetriever enables precise error locator and effective correction, enhancing performance for complex cases.

The differential impact of ContextRetriever between the NACA0012 Case 1 and Nozzle Case 5 reflects their distinct flow physics. The NACA0012 case, an incompressible flow simulation using the pisoFoam solver, exhibits weak pressure-velocity coupling, reducing the need for advanced error handling. Conversely, Nozzle Case 5, a compressible flow scenario with the rhoCentralFoam solver, involves strong interdependencies among density, temperature, and pressure, governed by the equation of state. The ContextRetriever module's advanced analysis of inter-file dependencies proves particularly effective for the Nozzle case's coupled physics, yielding greater performance gains compared to the less coupled NACA0012 case. This highlights ContextRetriever 's proficiency in managing complex physical interactions across diverse CFD scenarios.

Advancing to Configuration E, which integrates comprehensive article interpretation in Stage 2, maintained the 60% 10-step operational success rate for Nozzle Case 5 while increasing the accurate configuration success rate from 0% to 30%. This metric reflects precise adherence to source specifications. Consistent with findings for the NACA0012 case, Configuration E's enhanced setup extraction and validation ensure greater configuration fidelity across flow regimes, reinforcing ChatCFD's robustness in literature-derived cases.

The performance disparity between the literature-derived NACA0012 Case 1 and Nozzle Case 5 arises from two key factors: (1) the increased complexity of managing numerous files with intricate content in compressible flows, and (2) the stringent requirement for consistency across coupled thermophysical models (e.g.,

Figure 12: ChatCFD's performance metrics for accurate configuration of literature-derived NACA0012 Case 1 and Nozzle Case 5. Metrics include: number of reflection iterations (#reflection), total token consumption for DeepSeek-R1 and DeepSeek-V3 models (R1/V3 total token), execution cost (Cost $), and average token consumption per reflection round for both models (R1/V3 token/#reflection).

<!-- image -->

0/p , 0/T , constant/thermophysicalProperties ). The Nozzle case, a compressible flow simulation using rhoCentralFoam , involves strong coupling among pressure, density, temperature, and the equation of state, necessitating precise dimensional and physical alignment across governing equations and material properties. This contrasts with the NACA0012 case, an incompressible flow simulation using pisoFoam with weaker pressure-velocity coupling, which imposes fewer constraints. Additionally, an implicit bias in LLM training data, predominantly focused on incompressible flow scenarios, leads to erroneous defaults in compressible cases. For instance, in OpenFOAM, incompressible flows use kinematic pressure ( p k , dimensions L 2 T -2 )-dynamic pressure divided by density-while compressible flows require absolute thermodynamic pressure ( p , dimensions ML -1 T -2 ). The DeepSeek-R1 model's tendency to apply incompressible conventions to compressible cases introduces dimensional errors, such as incorrect pressure settings, requiring additional reflection iterations for correction.

Figure 12 quantifies these challenges: the Nozzle case requires over twice the reflection iterations, triple the DeepSeek-R1 tokens, and six times the DeepSeek-V3 tokens compared to the NACA0012 case, resulting in approximately triple the execution cost. The disproportionate DeepSeek-V3 token increase (6x vs. 3x for DeepSeek-R1 ) reflects frequent dimensional inconsistencies and persistent errors, necessitating more invocations of the costlier DeepSeek-V3 model. Figure 12(c) confirms this, showing a DeepSeek-V3 token multiplier of 2.7x per reflection for the Nozzle case, compared to 2x for DeepSeek-R1 , relative to the NACA0012 case. These metrics highlight the computational burden of complex flow simulations and underscore ChatCFD's ability through advanced error correction.

## 3.5.2 Influence of Turbulence Models on ChatCFD Efficacy

Figure 13 presents ChatCFD's performance across four turbulence models-Spalart-Allmaras (SA), kϵ , kω SST, and RNG kϵ -applied to the literature-derived NACA0012 cases (Cases 1 to 4). The analysis focuses on simulations achieving accurate configuration per published specifications, excluding those with only basic operational success (e.g., 10-step operational success rate). All results were obtained using the complete ChatCFD system (Configuration E).

As shown in Figure 13(a), success rates for accurate configuration vary significantly: SA at 40%, kω SST at 30%, kϵ at 20%, and RNG kϵ at 10%. This demonstrates ChatCFD's ability to handle diverse turbulence models, addressing LLM challenges noted in OpenFOAMGPT [31] through structured knowledge integration and RAG. Performance metrics in Figure 13(a,b) reveal consistent profiles for SA, kϵ , and kω SST, with an average of 6 reflection iterations, comparable DeepSeek-R1 token usage, and costs. In contrast, RNG kϵ requires 23.5 reflections-nearly 4x higher-elevating token consumption and costs.

Figure 13(c) attributes these disparities to configuration file token counts and model prevalence in OpenFOAM tutorials. While file sizes differ modestly ( 1,800 tokens for SA vs. 2,000 for others), tutorial distribution is more influential:

SA, kϵ , and kω SST appear in 28+ tutorials each, versus one for RNG kϵ . RNG kϵ 's performance deficit stems from: (1) its rarity in engineering applications, leading to sparse LLM training data and poor zero-shot configuration; and (2) limited ReferenceRetriever efficacy due to few OpenFOAM tutorials, hindering guidance for parameters in constant/turbulenceProperties , system/fvSchemes , and system/fvOption . This necessitates greater reliance on the LLM's limited intrinsic knowledge, reducing performance relative to common models.

Figure 13: ChatCFD's average performance for accurate configuration of NACA0012 Cases 1 to 4 across turbulence models (Spalart-Allmaras, kϵ , kω SST, RNG kϵ ). Metrics shown include: accuracy (success rate of precise configuration according to literature), number of reflection iterations (#reflection), total token consumption of the DeepSeek-R1 model (R1 total token), execution cost (Cost $), token count of configuration files (#token of config files), and the number of cases incorporating each turbulence model in OpenFOAM tutorials (#Case of model in OF tutorial).

<!-- image -->

## 3.6 Future Potential of ChatCFD in Multi-Agent CFD Workflows

Building on the insights from ChatCFD's performance in benchmark tutorial, perturbed variant, and literature-derived cases, ongoing work explores its potential as a foundational component in advanced engineering workflows. Beyond functioning as a standalone tool for CFD automation, ChatCFD is being adapted through the Model Context Protocol (MCP), a standardized framework for integrating specialized AI agents into a cohesive multi-agent system. MCP facilitates seamless communication and data exchange among agents by defining structured protocols for context sharing, enabling ChatCFD to collaborate with other intelligent tools to address complex, interdisciplinary engineering tasks.

As illustrated in Figure 14, this multi-agent workflow enables an end-to-end design cycle for CFD applications. The process begins with a 3D modeling agent that generates a three-dimensional geometry from a natural language description or a simple image, such as a sketch or photograph. A subsequent meshing agent automatically discretizes the geometry into a high-fidelity computational mesh, leveraging tools like commercial mesh generators or open-source alternatives (e.g., blockMesh , snappyHexMesh ). This mesh is then passed to the ChatCFD agent, which autonomously configures solver parameters, boundary conditions, and turbulence models (e.g., simpleFoam with Spalart-Allmaras or rhoCentralFoam with kω SST) to execute the CFD simulation, computing key performance metrics such as drag coefficients or pressure distributions. Finally, an optimization agent analyzes the simulation results, providing insights for iterative design refinement, such as adjusting geometry or boundary conditions to minimize drag or enhance flow stability. This MCP agent is available to use and can be found online (https://www.bohrium.com/apps/designagent0812).

This integrated approach highlights ChatCFD's transformative potential beyond simulation automation. By serving as a standardized computational module within a collaborative multi-agent ecosystem, ChatCFD can contribute to broader engineering design processes, from conceptual design to performance optimization. The MCP framework ensures interoperability, allowing ChatCFD to adapt to diverse tasks, such as aerodynamics, heat transfer, or multiphase flow

simulations, by interfacing with agents specialized in geometry generation, mesh refinement, or post-processing. This vision positions ChatCFD as a cornerstone for future AI-driven engineering workflows, paving the way for scalable, interdisciplinary applications in scientific discovery and industrial innovation.

Figure 14: ChatCFD integrated within a multi-agent workflow via the Model Context Protocol (MCP), enabling transformation from a single sentence or image to a complete CFD simulation.

<!-- image -->

## 4 Conclusion

This paper introduces ChatCFD, an innovative LLM-based CFD agent system for automating OpenFOAM simulations. By integrating a multi-agent architecture, domain-specific structured knowledge bases, precise error locator, and iterative reflection, ChatCFD dramatically outperforms prior systems while exposing fundamental limitations of current LLMs in scientific computing.

On 315 OpenFOAM benchmark cases, ChatCFD attains an 82.1% execution success rate -far surpassing MetaOpenFOAM (6.2%) and Foam-Agent (42.3%). More importantly, we introduced physical fidelity as a rigorous new metric: among all attempts, ChatCFD generates physically correct solutions in 68.12% of cases, with 60% of semantic failures traced to subtle boundary/initial condition errors that escape convergence checks. A dedicated Physics Interpreter further achieves 97.4% summary fidelity , revealing the striking gap between LLMs' fluency in high-level narration and their struggle to enforce dozens of tightly coupled domain constraints in executable code.

Resource-efficiency analysis on 205 benchmark cases further underscores ChatCFD's practical superiority. It consumes only 192.1k tokens and $0.208 per case -roughly half the tokens of Foam-Agent and 1.5 × cheaper than MetaOpenFOAM-while requiring the fewest reflection iterations on average. This efficiency gain over naive LLM prompting stems from strategic model routing (DeepSeek-V3 for generation, DeepSeek-R1 for reasoning) and a reflection memory that filters irrelevant context. Even under a strict iteration limit, ChatCFD's reflection trajectory consistently outperforms baselines, solving cases that saturate the others at 50% success after 25 iterations. These results demonstrate that structured domain integration not only boosts accuracy but dramatically reduces the operational cost of LLM-driven scientific computing, making high-fidelity CFD automation viable at scale.

Ablation studies confirm that structured OpenFOAM knowledge is indispensable: removing the Solver Template DB alone collapses accuracy to 48%, while the Error Locator Module proves the single most critical component for real-world robustness. Flexibility experiments demonstrate ChatCFD's ability to autonomously select appropriate solvers across compressible/incompressible and steady/transient regimes (95.23% success) and switch turbulence closures (100% success), even on unseen configurations.

By reproducing complex literature cases-NACA0012 airfoil, supersonic nozzles-ChatCFD shows 60-80% end-toend success where baselines fail entirely. These results establish potentially new benchmarks for AI-driven CFD and provide the quantitatively grounded evidence of where current LLMs break in scientific tool use.

ChatCFD effectively addresses the critical gaps in CFD automation agents. First, it integrates domain-specific knowledge through structured RAG and expert-designed agents, mitigating LLM training scarcity to enable precise boundary and initial condition specifications. Second, its modular, MCP-compatible design supports collaborative multi-agent networks, integrating specialized meshing agents to handle complex geometries beyond blockMesh . Third, by leveraging extensive literature-derived corpora with detailed specifications and benchmarks, it will facilitate

scalable testing and iterative agent evolution, harnessing millions of CFD papers for refinement and validation in the future. These advances will enable scientists across disciplines to explore bold ideas using CFD that were previously stalled by technical details. By liberating researchers from low-level case-setup burdens, it enables non-experts to perform computational fluid dynamics studies, allows rapid exploration of hypothesis and design spaces, and frees CFD specialists to focus on genuine scientific questions.

## ACKNOWLEDGEMENTS

This work is sponsored by the National Natural Science Foundation of China, Grant No. 92470127, and the Overseas Postdoctoral Talents Program in Guangdong Province.

## AUTHOR DECLARATIONS

## Conflict of Interest

The authors declare no conflicts of interest.

## Author Contributions

E Fan : Methodology, Data curation, Formal analysis, Writing - original draft. Kang Hu : Methodology, Software, Data curation, Formal analysis, Writing - original draft. Zhuowen Wu : Methodology. Jiangyang Ge : Data curation, Formal analysis. Jiawei Miao : Software, Resources. Yuzhi Zhang : Conceptualization, Software, Resources. He Sun : Resources, Writing - Review. Weizong Wang : Resources, Writing - Review, Funding acquisition. Tianhan Zhang : Conceptualization, Methodology, Formal analysis, Resources, Writing - Review &amp; Editing, Visualization, Supervision, Funding acquisition

## Data Availability

The data supporting the findings are available at https://github.com/ConMoo/ChatCFD .

## References

- [1] Jeffrey P Slotnick, Abdollah Khodadoust, Juan Alonso, David Darmofal, William Gropp, Elizabeth Lurie, and Dimitri J Mavriplis. Cfd vision 2030 study: a path to revolutionary computational aerosciences. Technical report, 2014.
- [2] E Fan, Weizong Wang, and Tianhan Zhang. Numerical investigation on flame dynamic and regime transitions during shock-cool flame interaction. Combustion and Flame , 273:113928, 2025.
- [3] Alfredo Iranzo. Cfd applications in energy engineering research and simulation: an introduction to published reviews. Processes , 7(12):883, 2019.
- [4] Joseph Amponsah, Emmanuel Adorkor, David Ohene Adjei Opoku, Anthony Ayine Apatika, and Vincent Nyanzu Kwofie. Computational modeling of hydrogen behavior and thermo-pressure dynamics for safety assessment in nuclear power plants. Physics of Fluids , 36(12), 2024.
- [5] Junjie Yao, Yuxiao Yi, Liangkai Hang, Weizong Wang, Yaoyu Zhang, Tianhan Zhang, Zhi-Qin John Xu, et al. Solving multiscale dynamical systems by deep learning. Computer Physics Communications , page 109802, 2025.
- [6] Zhengtong Li, Tingzhen Ming, Shurong Liu, Chong Peng, Renaud de Richter, Wei Li, Hao Zhang, and Chih-Yung Wen. Review on pollutant dispersion in urban areas-part a: Effects of mechanical factors and urban morphology. Building and Environment , 190:107534, 2021.
- [7] Yu-Hsuan Juan, Vita Ayu Aspriyanti, and Wan-Yi Chen. Wind flow characteristics in high-rise urban street canyons with skywalks. Physics of Fluids , 37(3), 2025.
- [8] Lian Shen, Yan Han, CS Cai, Peng Hu, Xu Lei, Pinhan Zhou, and Shuwen Deng. Equilibrium atmospheric boundary layer model for numerical simulation of urban wind environment. Physics of Fluids , 36(8), 2024.
- [9] Stefan Posch, Clemens Gößnitzer, Michael Lang, Ricardo Novella, Helfried Steiner, and Andreas Wimmer. Turbulent combustion modeling for internal combustion engine cfd: A review. Progress in Energy and Combustion Science , 106:101200, 2025.

- [10] Tianhan Zhang, Yuxiao Yi, Yifan Xu, Zhi X Chen, Yaoyu Zhang, Zhi-Qin John Xu, et al. A multi-scale sampling method for accurate and robust deep neural network to predict combustion chemical kinetics. Combustion and Flame , 245:112319, 2022.
- [11] Zhiwei Wang, Yaoyu Zhang, Pengxiao Lin, Enhan Zhao, Tianhan Zhang, Zhi-Qin John Xu, et al. Deep mechanism reduction (deepmr) method for fuel chemical kinetics. Combustion and Flame , 261:113286, 2024.
- [12] Xiaoyu Zhang, Yuxiao Yi, Lile Wang, Zhi-Qin John Xu, Tianhan Zhang, and Yao Zhou. Deep neural networks for modeling astrophysical nuclear reacting flows. The Astrophysical Journal , 2025.
- [13] Paul D Morris, Andrew Narracott, Hendrik von Tengg-Kobligk, Daniel Alejandro Silva Soto, Sarah Hsiao, Angela Lungu, Paul Evans, Neil W Bressloff, Patricia V Lawford, D Rodney Hose, et al. Computational fluid dynamics modelling in cardiovascular medicine. Heart , 102(1):18-28, 2016.
- [14] Siamak N Doost, Dhanjoo Ghista, Boyang Su, Liang Zhong, and Yosry S Morsi. Heart blood flow simulation: a perspective review. Biomedical engineering online , 15:1-28, 2016.
- [15] John David Anderson and John Wendt. Computational fluid dynamics , volume 206. Springer, 1995.
- [16] Jiri Blazek. Computational fluid dynamics: principles and applications . Butterworth-Heinemann, 2015.
- [17] Haixin Wang, Yadi Cao, Zijie Huang, Yuxuan Liu, Peiyan Hu, Xiao Luo, Zezheng Song, Wanjia Zhao, Jilin Liu, Jinan Sun, et al. Recent advances on machine learning for computational fluid dynamics: A survey. arXiv preprint arXiv:2408.12171 , 2024.
- [18] Jyotsna Balakrishna Kodman, Balbir Singh, and Manikandan Murugaiah. A comprehensive survey of open-source tools for computational fluid dynamics analyses. J. Adv. Res. Fluid Mech. Therm. Sci , 119:123-148, 2024.
- [19] Josh Achiam, Steven Adler, Sandhini Agarwal, Lama Ahmad, Ilge Akkaya, Florencia Leoni Aleman, Diogo Almeida, Janko Altenschmidt, Sam Altman, Shyamal Anadkat, et al. Gpt-4 technical report. arXiv preprint arXiv:2303.08774 , 2023.
- [20] Gemini Team, Rohan Anil, Sebastian Borgeaud, Jean-Baptiste Alayrac, Jiahui Yu, Radu Soricut, Johan Schalkwyk, Andrew M Dai, Anja Hauth, Katie Millican, et al. Gemini: a family of highly capable multimodal models. arXiv preprint arXiv:2312.11805 , 2023.
- [21] Daya Guo, Dejian Yang, Haowei Zhang, Junxiao Song, Ruoyu Zhang, Runxin Xu, Qihao Zhu, Shirong Ma, Peiyi Wang, Xiao Bi, et al. Deepseek-r1: Incentivizing reasoning capability in llms via reinforcement learning. arXiv preprint arXiv:2501.12948 , 2025.
- [22] Sirui Hong, Xiawu Zheng, Jonathan Chen, Yuheng Cheng, Jinlin Wang, Ceyao Zhang, Zili Wang, Steven Ka Shing Yau, Zijuan Lin, Liyang Zhou, et al. Metagpt: Meta programming for multi-agent collaborative framework. arXiv preprint arXiv:2308.00352 , 3(4):6, 2023.
- [23] Qingyun Wu, Gagan Bansal, Jieyu Zhang, Yiran Wu, Beibin Li, Erkang Zhu, Li Jiang, Xiaoyun Zhang, Shaokun Zhang, Jiale Liu, et al. Autogen: Enabling next-gen llm applications via multi-agent conversation. arXiv preprint arXiv:2308.08155 , 2023.
- [24] Jason Wei, Xuezhi Wang, Dale Schuurmans, Maarten Bosma, Fei Xia, Ed Chi, Quoc V Le, Denny Zhou, et al. Chain-of-thought prompting elicits reasoning in large language models. Advances in neural information processing systems , 35:24824-24837, 2022.
- [25] Shunyu Yao, Jeffrey Zhao, Dian Yu, Nan Du, Izhak Shafran, Karthik Narasimhan, and Yuan Cao. React: Synergizing reasoning and acting in language models. In International Conference on Learning Representations (ICLR) , 2023.
- [26] Noah Shinn, Federico Cassano, Edward Berman, Ashwin Gopinath, Karthik Narasimhan, and Shunyu Yao. Reflexion: Language agents with verbal reinforcement learning, 2023.
- [27] Timo Schick, Jane Dwivedi-Yu, Roberto Dessì, Roberta Raileanu, Maria Lomeli, Eric Hambro, Luke Zettlemoyer, Nicola Cancedda, and Thomas Scialom. Toolformer: Language models can teach themselves to use tools. Advances in Neural Information Processing Systems , 36:68539-68551, 2023.
- [28] Patrick Lewis, Ethan Perez, Aleksandra Piktus, Fabio Petroni, Vladimir Karpukhin, Naman Goyal, Heinrich Küttler, Mike Lewis, Wen-tau Yih, Tim Rocktäschel, et al. Retrieval-augmented generation for knowledgeintensive nlp tasks. Advances in neural information processing systems , 33:9459-9474, 2020.
- [29] Hrvoje Jasak, Aleksandar Jemcov, Zeljko Tukovic, et al. Openfoam: A c++ library for complex physics simulations. In International workshop on coupled methods in numerical dynamics , volume 1000, pages 1-20. Dubrovnik, Croatia), 2007.
- [30] Yuxuan Chen, Xu Zhu, Hua Zhou, and Zhuyin Ren. Metaopenfoam: an llm-based multi-agent framework for cfd. arXiv preprint arXiv:2407.21320 , 2024.

- [31] Sandeep Pandey, Ran Xu, Wenkang Wang, and Xu Chu. Openfoamgpt: A retrieval-augmented large language model (llm) agent for openfoam-based computational fluid dynamics. Physics of Fluids , 37(3), 2025.
- [32] Yue Ling, Somasekharan Nithin, Cao Yadi, and Pan Shaowu. Foam-agent: Towards automated intelligent cfd workflows. arXiv preprint arXiv:2505.04997 , 2025.
- [33] Zhehao Dong, Zhen Lu, and Yue Yang. Fine-tuning a large language model for automating computational fluid dynamics simulations. Theoretical and Applied Mechanics Letters , page 100594, 2025.
- [34] Edward J Hu, Yelong Shen, Phillip Wallis, Zeyuan Allen-Zhu, Yuanzhi Li, Shean Wang, Lu Wang, Weizhu Chen, et al. Lora: Low-rank adaptation of large language models. ICLR , 1(2):3, 2022.
- [35] Keyan Ding, Jing Yu, Junjie Huang, Yuchen Yang, Qiang Zhang, and Huajun Chen. Scitoolagent: a knowledgegraph-driven scientific agent for multitool integration. Nature Computational Science , pages 1-11, 2025.
- [36] VolcEngine. VolcEngine - A Cloud Service Platform by ByteDance. https://www.volcengine.com/ , 2025. Accessed: May 8, 2025.
- [37] Zhaozheng Sun and Wenhui Yan. Comparison of different turbulence models in numerical calculation of lowspeed flow around naca0012 airfoil. In Journal of Physics: Conference Series , volume 2569, page 012075. IOP Publishing, 2023.
- [38] T Yu, Y Yu, YP Mao, YL Yang, and SL Xu. Comparative study of openfoam solvers on separation pattern and separation pattern transition in overexpanded single expansion ramp nozzle. Journal of Applied Fluid Mechanics , 16(11):2249-2262, 2023.
- [39] Yuxuan Chen, Xu Zhu, Hua Zhou, and Zhuyin Ren. Metaopenfoam 2.0: Large language model driven chain of thought for automating cfd simulation and post-processing. arXiv preprint arXiv:2502.00498 , 2025.
- [40] Nigel Gregory and CL O'reilly. Low-speed aerodynamic characteristics of naca 0012 aerofoil section, including the effects of upper-surface roughness simulating hoar frost. 1970.

## A Appendix: detailed analysis of two literature-derived cases

Figure 15 presents a detailed analysis of the NACA0012 Case 1 performance across the five distinct Configurations (A through E), based on averaged results from 10-step operational success runs. This analysis yields several key insights into the system's operational efficiency and cost-effectiveness.

The evolution of reflection iterations, illustrated in Figure 15(a), demonstrates an encouraging improvement pattern. As retrieve modules were integrated from configuration A to C, we observed a substantial decrease in required reflection iterations, followed by stabilization. This observation aligns with MetaOpenFOAM's findings [39] regarding RAG positive impact on LLM response quality and CFD agent performance. The subsequent stability in reflection counts from configurations C to E suggests that additional RAG implementations, while beneficial for other aspects, had minimal impact on ChatCFD's reflection behavior. This plateau is reasonable given that ContextRetriever and comprehensive paper interpretation modules were primarily designed to enhance performance on physical-coupled cases and setup accuracy respectively, rather than affecting reflection behaviors for weak coupling cases.

An analysis of token consumption and associated costs, depicted in Figure 15(b), reveals a pattern of improving resource utilization efficiency. Transitioning from Configuration A to D, while the token consumption of the DeepSeek-R1 model remained relatively stable, the DeepSeek-V3 model's consumption decreased markedly. This reduction contributed to lower per-case operational costs (decreasing from $0.078 to $0.071). This cost improvement, concurrent with a significant increase in operational success rates from 20% to 80%, underscores the enhanced overall effectiveness of the ChatCFD system.

The transition from Configuration D to E resulted in a notable enhancement of system accuracy, as illustrated in Figure 11. This advancement was primarily achieved by improving the paper interpretation module's efficacy in extracting CFD case setups during Stage 2 and by ensuring greater consistency of these setups throughout the file correction processes in Stage 3. Although this upgrade led to increased DeepSeek-R1 token consumption and a corresponding 30% rise in operational cost (from $0.071 to $0.091), the substantial improvement in accurate case configuration-from 0% to 40%-justifies this investment.

Figure 15(c) illustrates the token consumption patterns per reflection iteration for both DeepSeek-R1 and DeepSeek-V3 models across the different Configurations. From Configuration A to E, the DeepSeek-R1 's consumption increases to more than two times, where the increasing consumption are mostly due to the integration of the ContextRetriever module in Configuration C. The DeepSeek-V3 's consumption initially rose due to expanded error handling capabilities and reference case integration through ReferenceRetriever . The subsequent stabilization of DeepSeek-V3 token

usage from Configuration C to E reflects the maturation of the reflection module's design, achieving optimal operational efficiency without necessitating additional DeepSeek-V3 module invocations.

Figure 15: Results for NACA0012 Case 1 using Configurations A to E, showing averages of 10-step operational success runs. Metrics shown include: number of reflection iterations (#reflection), total token consumption for DeepSeek-R1 and DeepSeek-V3 models (R1/V3 total token), execution cost (Cost $), and average token consumption per reflection round for both models (R1/V3 token/#reflection)

<!-- image -->

Figure 16 presents a comparative performance analysis for the Nozzle Case 5, specifically contrasting Configurations D and E, based on simulations that successfully completed ten steps. Configurations A through C were omitted from this comparison due to their inability to consistently reach this ten-step benchmark. The key differentiator between Configurations D and E is the integration of the comprehensive paper interpretation module's enhanced article interpretation capabilities in Configuration E. Notwithstanding this enhancement, Figure 16(a) indicates that the number of reflection iterations remained largely comparable between these two configurations, with only marginal reductions observed for Configuration E.

The deployment of comprehensive paper interpretation module resulted in increased token utilization, as depicted in Figure 16(b). Configuration E demonstrated higher token consumption by the DeepSeek-R1 model, whereas the DeepSeek-V3 model's token usage remained stable across both Configurations D and E. In terms of operational cost, the average per-case expenditure saw a modest increase from $0.28 to $0.31, an approximate 10% increment. Nevertheless, this additional investment yielded significant improvements in performance, particularly in case setup accuracy, which, as shown in Figure 11, rose substantially from 0% to 30%.

Further examination, illustrated in Figure 16(c), reveals that the average token consumption per reflection iteration for both the DeepSeek-R1 and DeepSeek-V3 models was comparable across Configurations D and E. This consistency underscores the efficacy of the system's architectural design, indicating that the comprehensive paper interpretation module detailed interpretation component operates as a distinct layer that augments overall system performance without adversely affecting the fundamental reflection mechanisms.

Figure 16: Results for Nozzle Case 5 using Configurations D and E, showing averages of 10-step operational success runs. Metrics shown include: number of reflection iterations (#reflection), total token consumption for DeepSeek-R1 and DeepSeek-V3 models (R1/V3 total token), execution cost (Cost $), and average token consumption per reflection round for both models (R1/V3 token/#reflection)

<!-- image -->

## B Appendix: CFD results of Cases 1 to 5

Figure 17: Velocity magnitude contours for NACA0012 Cases 1 to 4 at 10° angle of attack using four turbulence models. (a) Spalart-Allmaras model, (b) kϵ model, (c) kω SST model, (d) RNGkϵ model.

<!-- image -->

Figures 17 and 18 show the velocity magnitude contours and surface pressure coefficient distributions for the NACA0012 airfoil Cases 1 to 4 at 10° angle of attack using different turbulence models. As can be seen, the velocity contours from the four cases are similar, and the pressure coefficient distributions agree well with the experimental results [40].

Figure 19 illustrates the experimental and numerical schlieren images for the nozzle case at a nozzle pressure ratio of 3. The experimental results are reported in Ref. [38]. The numerical results exhibit a strong correspondence with the experimental observations, accurately capturing key flow features. The figure clearly depicts the formation of a Mach stem, resulting from the intersection and reflection of shock waves from the ramp and flap. This interaction also reveals the characteristic Mach reflection, distinguished by its λ -shock wave structure.

Table 2: Statistical averages of accurate configured cases for Naca0012 Case 1 and Nozzle Case 5.

| Metrics                                               | Naca0012, Case 1                                                                                                                         | Nozzle, Case 5                                                                                                                                                                         |
|-------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Physical models                                       | Flow type: Incompressible Solver: simpleFoam Turbulence: Spalart-Allmaras                                                                | Flow type: Compressible Solver: rhoCentralFoam Turbulence: Spalart-Allmaras Thermo model: hePsiThermo                                                                                  |
| Number of config- uration files                       | 9                                                                                                                                        | 12                                                                                                                                                                                     |
| Average tokens for all configura- tion files per case | 1,792                                                                                                                                    | 2,647                                                                                                                                                                                  |
| Configuration files                                   | 0/p 0/U 0/nut 0/nuTilda constant/transportProperties constant/turbulenceProperties system/controlDict system/fvSchemes system/fvSolution | 0/p 0/U 0/nut 0/nuTilda 0/T 0/alphat constant/transportProperties constant/turbulenceProperties constant/thermodynamicProperties system/controlDict system/fvSchemes system/fvSolution |
| DeepSeek-R1 call count                                | 19.1                                                                                                                                     | 48.3                                                                                                                                                                                   |
| DeepSeek-V3 call count                                | 17.4                                                                                                                                     | 35.2                                                                                                                                                                                   |

000000000

Figure 18: Comparison of pressure coefficients between CFD simulation and experimental data for NACA0012 airfoil at 10° angle of attack, using (a) Spalart-Allmaras model, (b) kϵ model, (c) kω SST model, and (d) RNGkϵ model. Symbols represent experimental results reported by Gregory et al. [40]

<!-- image -->

Figure 19: (a) Experimental schlieren image [38] and (b) numerical schlieren image of Case 5 for the Nozzle case at nozzle pressure ratio = 3.

<!-- image -->