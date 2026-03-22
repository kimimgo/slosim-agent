## Foam-Agent 2.0: An End-to-End Composable Multi-Agent Framework for Automating CFD Simulation in OpenFOAM

Ling Yue 1 † , Nithin Somasekharan 1 † , Tingwen Zhang 1 , Yadi Cao 2 , Shaowu Pan 1*

1* Rensselaer Polytechnic Institute.

2 University of California San Diego.

*Corresponding author(s). E-mail(s): pans2@rpi.edu;

†

These authors contributed equally to this work.

## Abstract

Computational Fluid Dynamics (CFD) is an essential simulation tool in engineering, yet its steep learning curve and complex manual setup create significant barriers. To address these challenges, we introduce Foam-Agent , a multi-agent framework that automates the entire end-to-end OpenFOAM workflow from a single natural language prompt. Our key innovations address critical gaps in existing systems: 1. An Comprehensive End-toEnd Simulation Automation: Foam-Agent is the first system to manage the full simulation pipeline, including advanced pre-processing with a versatile Meshing Agent capable of handling external mesh files and generating new geometries via Gmsh , automatic generation of HPC submission scripts, and post-simulation visualization via ParaView . 2. Composable Service Architecture: Going beyond a monolithic agent, the framework uses Model Context Protocol (MCP) to expose its core functions as discrete, callable tools. This allows for flexible integration and use by other agentic systems, such as Claude-code, for more exploratory workflows. 3. High-Fidelity Configuration Generation: We achieve superior accuracy through a Hierarchical Multi-Index RAG for precise context retrieval and a dependency-aware generation process that ensures configuration consistency. Evaluated on a benchmark of 110 simulation tasks, Foam-Agent achieves an 88.2% success rate with Claude 3.5 Sonnet, significantly outperforming existing frameworks (55.5% for Meta OpenFOAM ). Foam-Agent dramatically lowers the expertise barrier for CFD, demonstrating how specialized multi-agent systems can democratize complex scientific computing. The code is public at https://github.com/csml-rpi/Foam-Agent.

Keywords: Large Language Model Agents, Simulation Automation, AI4Science, Computational Fluid Dynamics

## 1 Introduction

The rapid progress in generative AI, particularly large language models (LLMs), has given rise to autonomous 'AI agents' capable of dramatically transforming workflows across many domains. These agents can break down complex tasks, invoke tools or other software, and iteratively refine their outputs with minimal human intervention. For example, AutoGPT [1] allows the user to plan and execute multi-step goals (such as researching a topic or writing code) by chaining LLM queries and actions autonomously. Numerous frameworks now enable LLMs to use external tools or models like Shen et al. [2] which leverage an LLM as a central controller to orchestrate a suite of HuggingFace models for solving multi-modal tasks. Frameworks such as ReAct [3] have focused on improving the decision-making capabilities of a single agent, interleaving reasoning and acting steps, allowing an LLM to generate plans, query tools or databases, and adjust its strategy based on intermediate results. Reflexion [4] introduced self-reflection and memory to enhance an agent's ability to correct its own mistakes, while Voyager [5] demonstrated open-ended exploration and tool use in an embodied setting like Minecraft. Toolformer [6] and Gorilla [7] extended the paradigm to reliable API calling and autonomous tool learning, ensuring agents can robustly interact with external systems. Multi-agent orchestration frameworks such as AutoGen [8] further show how multiple LLM-based agents can collaborate to solve complex problems. Multi-agent orchestration frameworks such as AutoGen [8] further show how multiple LLM-based agents can collaborate to solve complex problems. Beyond narrow task automation, LLM agents now demonstrate skills such as self-debugging code [9] and even autonomously conduct parts of scientific research (e.g. proposing hypotheses and analyzing data) [10]. In addition, generative agents , which simulate believable human-like behavior, have demonstrated how AI agents can handle complex interactive workflows in virtual environments [11]. Together, these advances suggest that AI agents can serve as general-purpose assistants, automating or accelerating tasks that traditionally required significant human effort.

In scientific research, AI-driven agents and models have become indispensable collaborators across domains like biology, chemistry, physics, and materials science. In structural biology, AlphaFold from DeepMind [12] achieved atomic-level accuracy in protein structure prediction, effectively solving a long-standing challenge in the field; more recently, AlphaFold 3 extended this capability to complexes comprising proteins, nucleic acids, and ligands [13]. In biomedical workflows, specialized LLM agents such as CRISPR-GPT [14] automate the design of gene-editing experiments-handling tasks like guide RNA selection with minimal human input, while systems like MedAgents [15] work through diagnostic reasoning given patient information, achieving zero-shot performance on medical reasoning tasks and Huang et al. [16], Jin et al. [17] proposing a general-purpose biomedical agent able to carry out a spectrum of research tasks. In chemistry, systems like ChemCrow [18] equip LLM agents with tool use to assist in synthesis planning, data querying, and reaction optimization, LLM-RDF [19] introduces an architecture of agents tailored for experimental chemistry tasks, including design, execution, and result interpretation and [20] introduces a multi-agent framework to execute quantum chemistry workflows from natural language input. In the field of materials science, several agentic frameworks have been introduced to combine available material science knowledge [21, 22] and provide tools for material science discovery. Further, Ghafarollahi and Buehler [23] showed that a multi-agent system can be used to discover novel high-performance alloys by coordinating tasks like phase diagram reasoning and a graph neural network to evaluate candidate compounds. These categories of autonomous systems reflect the rapidly expanding footprint of AI agents

Fig. 1 : Foam-Agent system architecture illustrating the complete end-to-end workflow from natural language input to post-processing visualization. The system features six primary agents with dynamic adaptive workflow topology: 1. Architect Agent : interprets user query and plans file and folder structures to be generated, 2. Meshing Agent : generates the OpenFOAM compatible mesh required to perform the simulation. Mesh can be produced by the agent itself using OpenFOAM native meshing modules or by using Gmsh meshing library. The user also has the option to provide the agent with externally developed mesh files in the form of .msh files or OpenFOAM native blockMeshDict / snappyHexMeshDict . 3. Input Writer Agent : generates OpenFOAM configuration files required to run the simulation like 0/U, 0/p, constant/physicalProperties, system/controlDict etc. 4. Runner Agent executes simulation either in the local environment of the user or in high performance computing environment as requested by the user in the prompt, 5. Reviewer Agent : diagnoses errors and proposes corrections through iterative debugging cycles. 6. Visualization Agent : generates visuals of physical quantities, if requested by the user within the user prompt, using Paraview/Pyvista python libraries. This adaptive workflow allows agents to autonomously determine execution paths based on intermediate results and problem complexity.

<!-- image -->

Fig. 2 : Modularized Foam-Agent architecture using model context protocol (MCP).

<!-- image -->

across scientific workflows-transforming not just individual tasks but entire research pipelines.

The application of agentic AI extends naturally from fundamental science discovery to the application engineering domains, which is characterized by well-established but laborious workflows and a heavy reliance on domain expertise of using sophisticated specialized tools. Agents with intelligent tool calling capabilities are being developed to automate tedious and repetitive tasks such as design modifications, data collection and analysis, visualization etc., thereby scaling engineering productivity and enabling the exploration of a wider design space. Frameworks such as AutoFEA [24] and MooseAgent [25], have demonstrated the ability to translate natural language descriptions of engineering problems into executable input files for solvers like simulation software, CalculiX [26], [27], and Abaqus [28]. These agents typically employ a step-by-step planning approach, decomposing the complex task of simulation setup into manageable sub-problems like geometry definition, meshing, material property assignment, boundary condition etc. Further, Tian and Zhang [29] demonstrated the importance of optimizing the role of agents and defining clear responsibilities to each agent. The development of benchmarks like FEABench [30] further drives progress in this area by providing a standardized means of evaluating an agent's ability to interact with FEA software APIs and solve real-world physics problems.

Similarly, in the field of Computational Fluid Dynamics (CFD), running a simulation is still a painstaking, expertise-heavy workflow. A practitioner must define the geometry, generate a mesh, specify boundary and initial conditions, choose turbulence and physics models, and configure solver parameters. Even small errors, such as a missing field definition or inconsistent units, can cause the solver to crash, making

debugging tedious and demanding deep domain knowledge. These barriers have historically limited accessibility and slowed innovation. Recent work shows that LLM-based agentic systems are beginning to transform CFD into a conversational, automated process. Studies like [31-33] introduces agents, which converts natural-language problem descriptions into complete, runnable OpenFOAM cases by retrieving templates and embedding domain knowledge. ChatCFD [34] offers an interactive, multi-modal conversational interface, enabling users to initialize and guide CFD simulations through. Together, these developments highlight a paradigm shift: CFD agents that reason, self-correct, and coordinate across the full pipeline, collapsing what was once days of manual setup into a single natural-language interaction.

Though these studies have shown initial promise but still face major shortcomings. First, they provide incomplete workflow coverage, often focusing narrowly on solver configuration while neglecting critical pre-processing of complex geometries and meshes, as well as essential post-processing tasks like visualization. Second, they struggle with high-fidelity generation for complex cases. Finally, most are designed as single, self-contained systems, which makes them difficult to integrate with other tools or agentic systems for broader, exploratory scientific workflows.

To address these challenges and advance the emerging paradigm of AI agents with specialized tool-usage capabilities, we introduce Foam-Agent (Figure 1), a multi-agent framework designed to interface with OpenFOAM: a complex scientific software used in CFD. We demonstrate its ability to fully automate OpenFOAM -based CFD workflows while also performing the pre-processing and post-processing tasks that are essential in scientific computing. FoamAgent treats CFD simulation as a dynamic sequence of tool mediated decisions: generating and integrating computational meshes, configuring solvers, diagnosing and recovering from errors, deploying to HPC environments, and producing validated post-processed outputs. In doing so, it demonstrates how LLM based agents can not only automate solver pipelines but also write tailored inputs specific to the tool being invoked, develop auxiliary scripts for post-processing and visualization of scientific data, and orchestrate heterogeneous software components into coherent workflows. Though Foam-Agent is demonstrated on automating CFD pipelines that use OpenFOAM , it can still be viewed as a template for extending agentic AI frameworks to complex domains of scientific computing, where adaptability, error recovery, and end-to-end integration are essential.

Similar to frameworks such as ChemCrow [35], SciToolAgent [36], and HoneyComb [21], Foam-Agent leverages large language models (LLMs) as central planners that invoke external tools to accomplish complex tasks from natural language prompts, integrating specialized software into the LLM's reasoning loop to handle domain-specific computations. It also parallels domain-focused scientific agents such as STELLA [17], Biomni [16], and AutoFEA [24] in tightly coupling agentic pipelines with complex scientific ecosystems. Across these diverse systems, a common paradigm emerges: lowering expertise barriers by enabling LLMs to reason about scientific problems, retrieve relevant knowledge, and delegate subtasks to computational tools.

Foam-Agent offers several key innovations. 1. Comprehensive End-to-End Simulation Automation : Foam-Agent is the first system to manage the full simulation pipeline, including advanced pre-processing with a versatile Meshing Agent capable of handling external mesh files and generating new geometries via Gmsh , automatic generation of HPC submission scripts, and post-simulation visualization via ParaView/Pyvista [37, 38] . 2. Composable Service Architecture : Instead of operating as a single, self-contained agent, the framework uses the Model Context Protocol (MCP) Figure 2 to expose its core functions as discrete, callable tools. This design allows for flexible integration and use by other agentic systems for more exploratory

workflows. 3. High-Fidelity Generation : We achieve superior accuracy through a Hierarchical Multi-Index RAG for precise context retrieval and a dependency-aware generation process that ensures configuration consistency across all files.

Extensive experiments on a comprehensive benchmark dataset of OpenFOAM cases, Foam-Agent achieves an 88.2% success rate with Claude 3.5 Sonnet, significantly outperforming existing approaches. Case studies are provided to demonstrate its ability to handle end to end workflow automation like generation/utilization of mesh files, visualization of the user specified flow fields, HPC job submission, integration of individual modules to generic agent orchestrator via MCP etc., which are the key novelties of our framework. This demonstrates the potential of specialized multi-agent systems to democratize access to complex scientific simulation tools. We public all code at https://github.com/csml-rpi/Foam-Agent.

## 2 Results

We evaluated Foam-Agent using a comprehensive benchmark dataset containing 110 OpenFOAM simulation cases across 11 distinct physics scenarios. The dataset spans a wide range of physical phenomena and geometric complexity Table 1, providing a thorough test of automated CFD capabilities.

Table 1 : Distribution of benchmark cases by physics category.

| Physics Category            |   Number of Cases |
|-----------------------------|-------------------|
| Shallow Water / Geophysical |                10 |
| Combustion / Reactive Flow  |                10 |
| Multiphase / Free Surface   |                10 |
| Shock Dynamics              |                10 |
| Turbulent Flow              |                20 |
| Laminar Flow                |                30 |
| Heat Transfer               |                20 |
| Total                       |               110 |

Each benchmark case is described using natural language prompts that include the problem description, physical scenario, geometry, solver requirements, boundary conditions, and simulation parameters. The success rate is measured by the percentage of cases that ran successfully through the agentic framework given the prompt describing the simulation scenarios.

Baseline Frameworks We compared Foam-Agent against two representative frameworks. MetaOpenFOAM [31], OpenFOAMGPT [39]. Since the authors of OpenFOAMGPT did not release their implementation, we used a variant of our own framework without the reviewer component to recreate their approach, which we refer to as OpenFOAMGPT-Alt throughout this paper. This implementation of our system (without the reviewer) closely mirrors their described functionality. For a comprehensive evaluation, each framework was tested with two frontier LLMs: Claude 3.5 Sonnet and GPT-4o.

Overall Performance Comparison Table 2 presents the comparative performance of the frameworks on executable success rate. Foam-Agent substantially outperforms both baselines across all tested models. With Claude 3.5 Sonnet, Foam-Agent achieves an 88.2% success rate compared to 55.5% for MetaOpenFOAM and 37.3% for

Table 2 : Comparison of executable success rates across frameworks.

| Base LLM Model    | Meta OpenFOAM   | OpenFOAM GPT-Alt   | Foam-Agent (Ours)   |
|-------------------|-----------------|--------------------|---------------------|
| Claude 3.5 Sonnet | 55.5%           | 37.3%              | 88.2%               |
| GPT-4o            | 17.3%           | 45.5%              | 59.1%               |

OpenFOAMGPT-Alt. With GPT-4o, Foam-Agent achieves 59.1% compared to 17.3% for MetaOpenFOAM and 45.5% for OpenFOAMGPT-Alt.

Ablation Studies Having established the superiority of Foam-Agent , we delve into the contribution of each component in our proposed framework. We analyze the impact of two key elements: the Reviewer Node and the File Dependency analysis module. All experiments were performed using the Claude-Sonnet-3.5 model. The results are summarized in Table 3.

Table 3 : Ablation study results evaluating the impact of the reviewer node and file dependency analysis on success rate, token usage, and reviewer efficiency. We varied the temperature to show that Foam-Agent 's performance is generalizable. The highest success rate for each configuration pair is highlighted in bold.

| Reviewer Node                                   | File Dependency                                 | Temp- erature                                   | Success Rate (%)                                | Token Usage                                     | Avg. Reviewer Loops                             |
|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|
| Configuration                                   | without Reviewer Node                           | without Reviewer Node                           | without Reviewer Node                           | without Reviewer Node                           | without Reviewer Node                           |
| ✗                                               | ✗                                               | 0.0                                             | 48.2                                            | 282,056                                         | N/A                                             |
| ✗                                               | ✓                                               | 0.0                                             | 56.4                                            | 314,670                                         | N/A                                             |
| ✗                                               | ✗                                               | 0.6                                             | 45.4                                            | 282,034                                         | N/A                                             |
| ✗                                               | ✓                                               | 0.6                                             | 57.3                                            | 314,922                                         | N/A                                             |
| Configuration with Reviewer Node (max loops=10) | Configuration with Reviewer Node (max loops=10) | Configuration with Reviewer Node (max loops=10) | Configuration with Reviewer Node (max loops=10) | Configuration with Reviewer Node (max loops=10) | Configuration with Reviewer Node (max loops=10) |
| ✓                                               | ✗                                               | 0.0                                             | 86.4                                            | 309,873                                         | 0.90                                            |
| ✓                                               | ✓                                               | 0.0                                             | 88.2                                            | 332,324                                         | 0.79                                            |
| ✓                                               | ✗                                               | 0.6                                             | 86.4                                            | 356,272                                         | 1.87                                            |
| ✓                                               | ✓                                               | 0.6                                             | 88.2                                            | 334,303                                         | 0.96                                            |

From the results in Table 3, it can be observed that the inclusion of the reviewer node is the most significant factor for performance. It dramatically improves the success rate from a baseline of roughly 50% to over 80% across all tested configurations. This highlights the critical role of iterative feedback and self-correction in solving complex scientific computing tasks. In the absence of reviewer, file dependency provides the most significant improvement on success rate, from 48.2% to 56.4% at T=0.0 and from 45.4% to 57.3% at T=0.6. However, because the reviewer operates independently of file dependency, the success rates do not differ significantly when the reviewer is present. As the workflow carries out more reviewer loops, any errors made during the initial file generation will be indiscriminately corrected by the reviewer. The evidence is the lower reviewer loops when file dependency is present (from 0.90 to 0.79 at T = 0.0 and from 1.87 to 0.96 at T=0.6). Therefore, the main application of file dependency is to help the reviewer converge, thus reducing API calls and workflow runtime.

The ablation study on RAG is summarized in Table 4. We first performed an experiment without reviewer to isolate the effect of hierarchical multi-index retrieval.

Table 4 : Comparsion of Foam-Agent's performance using the hierarchical multi-index (hierarchy) retrieval and a single-index (baseline) retrieval. All experiments are performed using Claude-Sonnet-3.5 at T=0.6. The highest success rates are highlighted in bold.

| Retrieval Method                                | Reviewer Node                                   | Success Rate (%)                                | Token Usage                                     | Avg. Reviewer Loops                             |                                                 |
|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|-------------------------------------------------|
| Configuration without Reviewer Node             | Configuration without Reviewer Node             | Configuration without Reviewer Node             | Configuration without Reviewer Node             | Configuration without Reviewer Node             | Configuration without Reviewer Node             |
| baseline                                        | ✗                                               | 44.6                                            | 271,136                                         | N/A                                             |                                                 |
| hierarchy                                       | ✗                                               | 57.3                                            | 306,844                                         | N/A                                             |                                                 |
| Configuration with Reviewer Node (max loops=10) | Configuration with Reviewer Node (max loops=10) | Configuration with Reviewer Node (max loops=10) | Configuration with Reviewer Node (max loops=10) | Configuration with Reviewer Node (max loops=10) | Configuration with Reviewer Node (max loops=10) |
| baseline                                        | ✓                                               | 84.6                                            | 323,011                                         | 0.73                                            |                                                 |
| hierarchy                                       | ✓                                               | 88.2                                            | 334,303                                         | 0.96                                            |                                                 |

The success rate of hierarchy (57.3%) is significantly higher than that of a singleindex retrieval (44.6%). Even with the reviewer, the effect of hierarchy is still noticible (88.2% vs. 84.6%). This multi-index approach significantly outperforms single-index approaches by reducing noise and improving retrieval precision. To address the semantic gap between natural language and technical terminology, we implement specialized tokenization and normalization procedures tailored to the CFD domain.

Case Studies To qualitatively assess the relative accuracy of Foam-Agent and Meta OpenFOAM , we visually compare their simulation outputs against the Ground Truth across representative benchmark cases ( CounterFlowFlame , wedge and forwardStep ) in Table 2. The CH 4 mass fraction distribution at the final timestep (0.5s) in CounterFlowFlame is shown in the top row, the middle row depicts the temperature distribution for the wedge case at a timestep of 0.2 and the bottom row shows the velocity magnitude distribution for forwardStep case at a timestep of 0.5s. The ground truth of CounterFlowFlame shows a sharp, well-defined concentration gradient characteristic of flame fronts. Meta OpenFOAM produces a more diffuse, inaccurate transition region with considerable deviations from the expected profile. Foam-Agent 's results closely match the ground truth. In the case of wedge , the ground truth exhibits a gradient of ve,velocity near the angled wall which is accurately captured by Foam-Agent , while Meta OpenFOAM fails to capture even the geometry of the problem. Finally in the case of forwardStep , Foam-Agent results are virtually indistinguishable from the ground truth, while Meta OpenFOAM results predict very low value of velocity within the domain. These cases demonstrate the accuracy of Foam-Agent in simulating canonical problems which are the stepping stones to more complicated flow fields.

External Mesh File Ability to process externally developed mesh files is a key novelty of our framework . We provide mesh in the form of .msh files to the framework. The boundary names and flow scenario is described in the prompt to the framework. We demonstrate this functionality on two cases: 1) Flow over multi-element airfoil [40]: This case describes the flow around multiple airfoils within the domain, with the flow of one affecting the flow of another. This 2D simulation uses an inlet velocity of 9 m/s and a fluid kinematic viscosity of 1 . 5 × 10 -5 m 2 / s. The simulation is set to use simpleFOAM solver and Spalart-Allmaras turbulence model, with a timestep of 1.0 s and a final time of 1000 s. The user prompt for the case is shown in Section B.1. 2) Flow over tandem wing configuration [40]: This case describes the flow around

<!-- image -->

Visualization generated by Foam-Agent

Fig. 3 : Comparison of simulation results produced by Meta OpenFOAM and Foam-Agent for CounterFlowFlame , wedge and forwardStep cases against the human expert generated ground truth. The top row shows the CH 4 mass fraction distribution comparison in CounterFlowFlame case at t=0.5s, the middle row shows the temperature distribution at a timestep of 0.2s for the wedge case and the bottom row shows the velocity magnitude distribution for forwardStep case at a timestep of 4s. Ground truth (left), Meta OpenFOAM (middle) and Foam-Agent (right). The visuals for Foam-Agent are auto generated by the agent. Foam-Agent generates a python script, which upon running loads the case and produces the visual as a .png file. The same python script is used to generate the visuals for ground truth and Meta OpenFOAM results.

a tandem configuration, with one wing located in the wake of the other. This 3D simulation uses an inlet velocity of 9 m/s and a fluid kinematic viscosity of 1 . 5 × 10 -5 m 2 / s. The simulation is set to use simpleFOAM solver and Spalart-Allmaras turbulence model, with a timestep of 1.0 s and a final time of 1000 s. The user prompt for the case is shown in Section B.2. The prompts along with the path to the mesh is provided to Foam-Agent for simulating the flow field. The corresponding mesh and the contour of the velocity magnitude at the final timestep for these cases is shown in Figure 4. Further we also compare Foam-Agent results with human expert generated simulation results. The simulation results from Foam-Agent closely matches with the result from human OpenFOAM expert as shown in the velocity magnitude contour plots. Inspecting the case files generated by the agent against those prepared by humans, we found them to closely match in key aspects such as physical parameters, turbulence model selection, boundary condition assignment, and solver choice. Minor differences were observed in the numerical schemes, but these did not significantly affect the outcomes, as confirmed by the velocity contours.

<!-- image -->

↑ Foam-Agent accepts path to the mesh file from user

Fig. 4 : External mesh file processing ability of Foam-Agent analyzed through a 2D multi-element airfoil (top) and tandem wing case (bottom). The prompt describing the case along with path of the .msh file is given to Foam-Agent . The prompt also mentions that Foam-Agent needs to generate visuals at the final timestep for the velocity magnitude Section B.1 and Section B.2. The architect agent plans the required files needed for the simulation. The external mesh is converted to OpenFOAM compatible format by the meshing agent. With other agentic nodes handling the generation of input files, execution, error correction etc. At the end of workflow, the visualization agent generates a python file, which is also run by the agent to create visuals of the velocity magnitude. The same python file is used in generating the visuals of the human expert generated result to maintain consistency of the qualitative visualization.

Gmsh Based Mesh Generation We demonstrate the mesh generation capabilities of Foam-Agent utilizing Gmsh python library using two cases, where the natural language description of the mesh is provided to the agent. 1) Flow over cylinder : The description of the geometry given to the agent is detailed in Section B.3. 2) Flow over two square obstacles : The description of the geometry given to the agent is detailed in Section B.4. The simulation is run for 10s using pisoFOAM solver. The mesh generated by the agent and the contour of velocity magnitude at the final timestep for the two cases is shown in Figure 5. To underscore the necessity of a specialized meshing agent that leverages external tools such as Gmsh , we compare meshes generated with OpenFOAM 's native meshing utilities against those produced via the Gmsh Python API. The native meshing modules were unable to capture the intended geometry of the flow scenario, whereas the Gmsh-based approach successfully generated the overall domain and accurately represented the obstacles within it. The meshing agent creates the mesh in the form of .msh file which is then converted to OpenFOAM compatible format and then further used for simulating the flow field.

HPC Runner We demonstrate the capabilities of the HPC Runner Agent by instructing the framework to perform a 3D lid-driven cavity flow with one million cells, following the setup used in [41], using the prompt shown in Section B.5. The agent generates the necessary OpenFOAM case files along with a Slurm submission script for the HPC platform Perlmutter . The platform name and submission account number are provided in the prompt, while the agent relies on the LLM's knowledge of the system's documentation to produce correct Slurm syntax, which can vary across clusters. The resulting highly refined mesh, the velocity contours at the final timestep (midz slice)

Fig. 5 : Flow over a cylinder (top) and flow over two square obstacle (bottom) simulated by Foam-Agent using gmsh meshing library. The user prompt describes that the mesh is to be made using Gmsh and the meshing agent generates the python file required to generate the mesh file utilizing the Gmsh library in python. The execution of the python creates the mesh in .msh format, which is then converted to OpenFOAM compatible format. The flow then proceeds to the other agents handling the generation of input files, running, error correction etc. Finally the agent creates the visuals of the velocity magnitude at the final magnitude as requested by the user within the prompt (Section B.3 and Section B.4) to Foam-Agent . To emphasize the necessity of Gmsh based mesh generation within the meshing agent, we show the mesh generated by OpenFOAM's native meshing modules given the description of the flow domain. The native meshing modules are unable to capture the obstacles in the domain and produces incorrect simulation result as shown the visualization of the velocity magnitude.

<!-- image -->

and corresponding streamlines are presented in Figure 6. The automatically generated Slurm script is given in Section C.

## 3 Methods

## 3.1 System Architecture Overview

Foam-Agent employs a modular architecture comprising six primary components. This multi-agent framework interprets the natural language requirement from the

<!-- image -->

subdomains.

Fig. 6 : Highly refined simulation of a 3D lid driven cavity flow. Due to the number of cells used in the simulation the running of this case requires high performance computing (HPC) resources. The prompt suggest the cluster on which the simulation is to be run and the number of cores over which the domain is to be decomposed. The agent then creates all the required OpenFOAM files for the simulation along with the slurm script needed to submit the job in HPC. The mesh is divided into 32 subdomains as mentioned in the prompt and each domain is assigned a core on which the computations for that sub-domain is carried out. Further the 2D slice of the flow field at the final timestep generated by the agent and the 3D manually generated streamlines are also shown here.

user, leverages Retrieval-Augmented Generation (RAG) to find similar cases and configurations within tutorial database, generates simulation configurations, execute simulations, analyze errors, and implement corrections through iterative refinement. The six primary components form the backbone of the system are Architect Agent , Meshing Agent , Input Writer Agent , Runner Agent , Reviewer Agent , Visualization Agent .

## 3.2 Agent Components

Architect Agent The Architect Agent translates natural language descriptions into structured simulation plans through a three-stage process. First, in the Requirement Classification stage, it analyzes the user requirement to classify the simulation according to domain-specific taxonomies, employing Pydantic models for structured output validation. Second, during Reference Case Retrieval, it queries the hierarchical indices to identify semantically similar cases, using a cascading approach that refines initial matches with detailed structural information. Finally, in the Simulation Decomposition stage, it decomposes the task into required files and directories, creating a detailed plan specifying file dependencies and generation priorities. The output for a simulation requirement R is a structured plan P = { F 1 , F 2 , ..., F n } where each F i represents a required file with its dependencies, content requirements, and generation priority.

Meshing Agent There multiple ways that Foam-Agent handles the meshing (preprocessing) stage of the simulation: 1) OpenFOAM native: The agent generated the required OpenFOAM native files for mesh generation like blockMeshDict and/or snappyHexMeshDict . It then continues to generate the polyMesh folder which contains the detailed mesh information for OpenFOAM to process the case by running

commands such as blockMesh and/or snappyHexMesh . The agent has the complete autonomy in this scenario. 2) External Mesh Files: Foam-Agent does not generate the mesh files for the simualtion. The user has the ability to provide the agent with mesh specific files either in the form of native OpenFOAM dictionaries ( blockMeshDict , snappyHexMeshDict ) or pre-generated meshes from external tools (e.g., .msh ). Given any of these inputs, the agent converts the mesh to OpenFOAM format, generating the ployMesh folder. 3) Gmsh: Foam-Agent can take in natural language input of the physical domain, the boundary names and generate a mesh file (.msh) based on this description. The agent creates a python script representing the geometry and mesh using Gmsh python library and with further execution generating the mesh file. This mesh file is then converted to OpenFOAM format. This functionality was found be to required as OpenFOAM 's inherent meshing creation capabilities were found to be insufficient in creating novel geometries not present in the tutorials. The agent intelligently chooses one of these options based on the user requirements.

Input Writer Agent Following established human CFD expertise, the Input Writer Agent implements a structured file generation sequence that respects OpenFOAM 's hierarchical organization: it begins with the system directory (simulation control parameters and numerical schemes), proceeds to the constant directory (physical, turbulence properties), then the 0 directory (initial and boundary conditions), and finally produces auxiliary files (e.g., Allrun ) for executing the simulation. This ordering naturally enforces dependencies, as files prescribing boundary conditions depend on turbulence model/physical properties provided in the constant directory, which themselves can be dependent on solver configurations. Formally, the process can be represented as a directed acyclic graph G = ( V, E ), where vertices V denote files and edges E capture their dependencies. To maintain internal consistency, the agent employs Contextual Generation , where the context C i of each file F i includes all previously generated files, enabling coherent parameter referencing across the configuration. Finally, schema validation is applied through Pydantic models, ensuring that generated files satisfy both the syntactic and semantic constraints of OpenFOAM , thereby reducing generation errors.

Runner Agent The Runner Agent interfaces with the OpenFOAM execution environment by preparing the simulation (cleaning artifacts, setting up output capture) and running the Allrun script. Simulations can be executed either on the user's local machine or deployed to HPC clusters , as specified in the requirements. For HPC runs, the agent automatically generates Slurm scripts, submits jobs with the provided account number, and monitors their progress until completion. It can use promptspecified parameters (e.g., nodes, processes per node) or infer them from the mesh decomposition and problem size. This capability enables Foam-Agent to scale seamlessly from desktop prototyping to large-scale industrial CFD workloads. At the end of simulation, it performs critical error detection by analyzing these logs to identify specific error patterns, extracting relevant messages and contextual information for subsequent analysis by other components. The error detection process is formalized as a pattern matching function E : L →{ e 1 , e 2 , ..., e m } that maps execution logs L to a set of structured error records e i , each containing the error message, location, and severity.

Reviewer Agent The Reviewer Agent implements error analysis and correction of the case being run. It performs Error Contextualization by assembling detailed context including error messages, affected files, and reference cases. It then conducts Review Trajectory Analysis by maintaining a record of previous attempts as H = { ( F 1 i , E 1 ) , ( F 2 i , E 2 ) , . . . , ( F n i , E n ) } , where F j i represents the file state in iteration

j and E j represents the resulting errors. Finally, it executes Solution Generation by producing structured modifications to problematic files while maintaining consistency with other configuration elements. The correction process can formalized as an optimization problem: finding the minimal set of changes ∆ F to files F that eliminates errors E while respecting user constraints C and maintaining consistency across all files (Algorithm 1). This review-correction cycle is repeated iteratively until either the errors are resolved or a maximum number of attempts specified by the user is reached (Figure 1).

Visualization Agent After the runner agent produces a successful run, Foam-Agent verifies if visualization is requested by the user within the provided prompt. If yes, the agent will try to understand the quantity to be visualized from the prompt and generate a python script utilizing either Pyvista library or paraview python routines (the user can choose) to generate the required visualization. The agent will then execute the python file and correct errors if any (like the reviewer agent) till the visualization(s) is saved as a .png file in the run directory. The maximum number of tries can be set by the user.

## Algorithm 1 Iterative Refinement Process

```
Require: Natural language requirement R , maximum iterations M Ensure: Final simulation configuration F ∗ 1: P ← Architect( R ) ▷ Generate initial plan 2: G ← Meshing( R ) ▷ Generate Computational Mesh 3: F 0 ← InputWriter( P ) ▷ Initial file generation 4: H ←{} ▷ Initialize history 5: for i ← 1 to M do 6: L i , S i ← Runner( F i -1 ) ▷ Execute simulation 7: if S i = SUCCESS then 8: return F i -1 ▷ Simulation successful 9: end if 10: H ← H ∪ { ( F i -1 , L i ) } ▷ Update history 11: E i ← ParseErrors( L i ) ▷ Extract errors 12: ∆ F i ← Reviewer( E i , F i -1 , H ) ▷ Generate corrections 13: F i ← Apply( F i -1 , ∆ F i ) ▷ Apply corrections 14: end for 15: return F M ▷ Return best attempt
```

## 3.3 RAG

A key innovation in Foam-Agent is its hierarchical multi-index retrieval system (Algorithm 2) that segments domain knowledge into specialized indices optimized for specific phases of the simulation workflow. This approach significantly improves retrieval precision compared to conventional single-index RAG systems.

Knowledge Base Organization We construct the knowledge base by parsing OpenFOAM 's tutorial cases, extracting information across four dimensions. The first dimension is Case Metadata , which includes fundamental attributes such as case name, flow domain, physical category, and solver selection. The second dimension is Directory Structures , which captures the hierarchical organization of files and directories in reference cases. The third dimension is File Contents , which preserves configuration file content, including syntax, parameter definitions, and commenting.

The fourth dimension is Execution Scripts , which includes command sequences for preparation, execution, and post-processing.

Specialized Vector Index Architecture Rather than using a monolithic database, Foam-Agent implements four distinct FAISS [42] indices, each serving a specific purpose. The Tutorial Structure Index encodes high-level case organization patterns for identifying appropriate structural templates. The Tutorial Details Index contains configuration details for boundary conditions, numerical schemes, and physical models. The Execution Scripts Index stores execution workflows for generating appropriate command sequences. The Command Documentation Index maintains utility documentation for correct command usage and parameter selection. Each index employs a 1536-dimensional vector with text-embedding-3-small model from OpenAI. The retrieval process is given in Algorithm 2.

## Algorithm 2 Hierarchical Multi-Index Retrieval

| Require: Query q , stage s , previous context Ensure: Retrieved context r 1: E ← Embed( q ) 2: I ← SelectIndex( s ) 3: R i ← TopK( I,E,k = 5) 4: R f ← FilterByRelevance( R i , q, c ) 5: r ← FormatContext( R f , s ) 6: return r   | ▷ Generate query embedding ▷ Select appropriate index for stage ▷ Retrieve top-k matches ▷ Filter by relevance ▷ Format based on stage   |
|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------|

## 3.4 Decoupling Capabilities via a Model Context Protocol (MCP)

To transition Foam-Agent from a monolithic tool into a flexible scientific service, we design its core around the Model Context Protocol (MCP). This decouples the CFD workflow into atomic, callable functions exposed via a standardized protocol, making Foam-Agent a composable component that higher-level agents or workflow engines can orchestrate. The MCP design follows three principles: atomicity (each function does one task), statefulness (tracking multi-stage simulations via identifiers such as case id and job id etc.), and workflow decoupling (separating meshing, solving, and post-processing). These features maximize flexibility while preserving fine-grained control, with the key functions summarized in Table 5.

While the MCP provides the foundational capabilities, a robust framework is required to orchestrate these functions. We achieve this through the following.

Stateful Workflow Orchestration with LangGraph The sequence of MCP function calls is not fixed; it is dynamically determined by an intelligent orchestrator. We implement this orchestrator as a stateful graph using LangGraph. The nodes in the graph correspond to calls to the MCP functions, while the edges represent conditional logic that directs the workflow based on the outcomes of each step.

Ensuring Reliability with Structured I/O A key challenge in LLM-tool interaction is avoiding errors from malformed or inconsistent data. To ensure reliability, we enforce strict schemas for all data exchanges, including MCP function I/O and LangGraph state. Using Pydantic, we define explicit data models that enable runtime validation and type checking, establishing a clear contract between the LLM, orchestrator, and tool functions.

Table 5 : The core functions of the Foam-Agent Model Context Protocol (MCP). Each function represents a decoupled capability within the CFD workflow, featuring strongly-typed inputs and outputs to ensure reliable interaction with orchestrating agents.

| Function Name             | Description                                                                                              | Input Schema                     | Output Schema                      |
|---------------------------|----------------------------------------------------------------------------------------------------------|----------------------------------|------------------------------------|
| create case               | Initializes a new CFD simulation case and its workspace.                                                 | { user prompt: str }             | { case id: str }                   |
| plan simulation structure | (Architect Agent) Plans the required file and directory structure based on the user prompt.              | { case id: str }                 | { plan: List[ { file, folder } ] } |
| generate file content     | (Input Writer Agent) Generates the content for a single specified configuration file.                    | { case id, file, folder }        | { content: str }                   |
| generate mesh             | (Meshing Agent) Asynchronously generates the computational mesh using a specified method.                | { case id, mesh config: Dict }   | { job id: str }                    |
| generate hpc script       | (HPC Agent) Generates a job submission script (e.g., Slurm) for a high-performance computing cluster.    | { case id, hpc config: Dict }    | { script content: str }            |
| run simulation            | (Runner Agent) Asynchronously executes the simulation either locally or by submitting to an HPC cluster. | { case id, environment: str }    | { job id: str }                    |
| check job status          | Checks the status of any asynchronous job (meshing, simulation, visualization).                          | { job id: str }                  | { status: Dict }                   |
| get simulation logs       | Retrieves detailed logs for a failed job to enable error diagnosis.                                      | { case id, job id }              | { logs: Dict }                     |
| review and suggest fix    | (Reviewer Agent) Analyzes error logs and proposes corrective actions.                                    | { case id, logs }                | { suggestions: Dict }              |
| apply fix                 | Applies suggested modifications to the relevant case files.                                              | { case id, modifications: List } | { status: str }                    |
| generate visualization    | (Visualization Agent) Asynchronously generates a visualization of the simulation results.                | { case id, quantity, ... }       | { job id: str }                    |

Achieving Observability with LangSmith The complex and often nondeterministic behavior of multi-agent systems makes them notoriously difficult to debug and analyze. To address this, we integrate LangSmith for end-to-end traceability. Every MCP function call, LLM invocation, and state transition within the LangGraph orchestrator is logged into an immutable record of the agent's 'thought process' and actions. This allows precise reconstruction of a simulation's execution,

including intermediate inputs, outputs, and errors, providing the observability needed for debugging, performance analysis, and ensuring scientific reliability.

## 4 Conclusion

This work introduced Foam-Agent , a modular multi-agent framework that automates end-to-end CFD workflows using OpenFOAM from natural language prompts. Our benchmark evaluation across 110 diverse cases demonstrated an 88.2% success rate, substantially surpassing prior frameworks. Extending beyond its core capabilities, Foam-Agent introduces integrated pre-processing, allowing users to import external mesh files or generate meshes on demand using external tools such as Gmsh , while also providing seamless execution on high-performance computing (HPC) platforms. Its decoupled architecture, built on the Model Context Protocol and LangGraph orchestration, enables seamless integration with other agentic ecosystems. All these capabilities not only lower the barrier of entry for non-experts but also empower experienced engineers to accelerate complex workflows. Foam-Agent establishes a foundation for scalable and trustworthy scientific automation. Future directions include extending its modular services to other solvers, enhancing physics coverage. Ultimately, this framework represents a step toward democratizing access to high-fidelity simulation and shaping the role of intelligent agents in computational science.

Acknowledgements. This work was supported by U.S. Department of Energy under Advancements in Artificial Intelligence for Science with award ID DESC0025425. The authors thank the Center for Computational Innovations (CCI) at Rensselaer Polytechnic Institute (RPI) for providing computational resources during the early stages of this research. Numerical experiments are performed using computational resources granted by NSF-ACCESS for the project PHY240112 and that of the National Energy Research Scientific Computing Center, a DOE Office of Science User Facility using the NERSC award NERSC DDR-ERCAP0030714. S. Pan is supported by Google Research Scholar Program.

## References

- [1] Yang, H., Yue, S., He, Y.: Auto-GPT for online decision making: Benchmarks and additional opinions. arXiv preprint arXiv:2306.02224 (2023)
- [2] Shen, Y., Song, K., Tan, X., Li, D., Lu, W., Zhuang, Y.: HuggingGPT: Solving AI tasks with ChatGPT and its friends in Hugging Face. In: Proc. of the 37th Conference on Neural Information Processing Systems (NeurIPS) (2023)
- [3] Yao, S., Zhao, J., Yu, D., Du, N., Shafran, I., Narasimhan, K., Cao, Y.: ReAct: Synergizing Reasoning and Acting in Language Models. Preprint at https://arxiv. org/abs/2210.03629 (2023)
- [4] Shinn, N., Labash, B., Gopinath, A.: Reflexion: an autonomous agent with dynamic memory and self-reflection. arXiv preprint arXiv:2303.11366 2303 (11366) (2023)
- [5] Wang, G., Qin, Y., Kosaraju, V., Lee, D., Zhang, F., Liang, P., Chen, J., Chen, Z., Ilievski, I., et al.: Voyager: an open-ended embodied agent in minecraft. arXiv preprint arXiv:2305.16291 2305 (16291) (2023)
- [6] Schick, T., Dwivedi-Yu, J., Dess` ı, R., Raileanu, R., Lomeli, M., Zettlemoyer, L., Cancedda, N., Scialom, T.: Toolformer: Language models can teach themselves to use tools. arXiv preprint arXiv:2302.04761 2302 (04761) (2023)
- [7] Patil, S.G., Zhang, T., Wang, X., Gonzalez, J.E.: Gorilla: Large language model connected with massive apis. In: Advances in Neural Information Processing Systems (NeurIPS) (2024). Preprint at urlhttps://arxiv.org/abs/2305.15334
- [8] Wu, T., Zhang, S., Zheng, A., et al.: Autogen: Enabling next-gen llm applications via multi-agent conversation. arXiv preprint arXiv:2308.08155 2308 (08155) (2023)
- [9] Dong, Y., Jiang, X., Jin, Z., Li, G.: Self-collaboration code generation via ChatGPT. ACM Trans. Softw. Eng. Methodol. 33 (7), 1-38 (2024)
- [10] Boiko, D.A., MacKnight, R., Gomes, G.: Emergent autonomous scientific research capabilities of large language models. arXiv preprint arXiv:2304.05332 (2023)
- [11] Park, J.S., O'Brien, J., Cai, C., Morris, M.R., Liang, P., Bernstein, M.: Generative agents: Interactive simulacra of human behavior. In: Proc. of the 36th ACM Symposium on User Interface Software and Technology (UIST), pp. 1-22 (2023)
- [12] Jumper, J., et al. : Highly accurate protein structure prediction with alphafold. Nature 596 (7873), 583-589 (2021)
- [13] Abramson, J., et al.: Accurate structure prediction of biomolecular interactions with alphafold 3. Nature (2024)
- [14] Huang, K., Qu, Y., Cousins, H., Johnson, W.A., Yin, D., Shah, M., Zhou, D., Altman, R., Wang, M., Cong, L.: CRISPR-GPT: An LLM agent for automated design of gene-editing experiments. arXiv preprint arXiv:2404.18021 (2024)

- [15] Tang, X., Zou, A., Zhang, Z., Li, Z., Zhao, Y., Zhang, X., Cohan, A., Gerstein, M.: MedAgents: Large language models as collaborators for zero-shot medical reasoning. arXiv preprint arXiv:2311.10537 (2023)
- [16] Huang, K., Zhang, S., Wang, H., Qu, Y., Lu, Y., Roohani, Y., Li, R., Qiu, L., Li, G., Zhang, J., Yin, D., Marwaha, S., Carter, J.N., Zhou, X., Wheeler, M., Bernstein, J.A., Wang, M., He, P., Zhou, J., Snyder, M., Cong, L., Regev, A., Leskovec, J.: Biomni: A general-purpose biomedical ai agent. bioRxiv (2025) https://doi.org/10.1101/2025.05.30.656746 https://www.biorxiv.org/content/early/2025/06/02/2025.05.30.656746.full.pdf
- [17] Jin, R., Zhang, Z., Wang, M., Cong, L.: STELLA: Self-Evolving LLM Agent for Biomedical Research (2025). https://arxiv.org/abs/2507.02004
- [18] Bran, A.M., et al. : Chemcrow: augmenting large-language models with chemistry tools. Nature Machine Intelligence 6 (4), 525-535 (2024)
- [19] Anonymous: An automatic end-to-end chemical synthesis development platform based on llm-rdf. Nature Communications (2024)
- [20] Zou, Y., Cheng, A.H., Aldossary, A., Bai, J., Leong, S.X., Campos-GonzalezAngulo, J.A., Choi, C., Ser, C.T., Tom, G., Wang, A., Zhang, Z., Yakavets, I., Hao, H., Crebolder, C., Bernales, V., Aspuru-Guzik, A.: El agente: An autonomous agent for quantum chemistry. Matter 8 (7), 102263 (2025) https://doi.org/10. 1016/j.matt.2025.102263
- [21] Hou, Z., et al.: Honeycomb: a flexible llm-based agent system for materials science. arXiv preprint 2409 (00135) (2024)
- [22] Buehler, M.J.: MechGPT: A language-based strategy for mechanics and materials modeling that connects knowledge across scales, disciplines, and modalities. Appl. Mech. Rev. 76 (2), 021001 (2024)
- [23] Ghafarollahi, A., Buehler, M.J.: Rapid and automated alloy design with graph neural network-powered LLM-driven multi-agent systems. arXiv preprint arXiv:2410.13768 (2024)
- [24] Hou, S., Johnson, R., Makhija, R., Chen, L., Ye, Y.: Autofea: enhancing ai copilot by integrating finite element analysis using large language models with graph neural networks. AAAI'25/IAAI'25/EAAI'25. AAAI Press, ??? (2025). https:// doi.org/10.1609/aaai.v39i22.34582 . https://doi.org/10.1609/aaai.v39i22.34582
- [25] Zhang, T., Liu, Z., Xin, Y., Jiao, Y.: MooseAgent: A LLM Based Multi-agent Framework for Automating Moose Simulation (2025). https://arxiv.org/abs/ 2504.08621
- [26] Dhondt, G., Wittig, K.: CalculiX: A Free Software Three-Dimensional Structural Finite Element Program. https://github.com/Dhondtguido/CalculiX. Accessed: 2025-09-03 (2025)
- [27] Permann, C.J., Gaston, D.R., Andrˇ s, D., Carlsen, R.W., Kong, F., Lindsay, A.D., Miller, J.M., Peterson, J.W., Slaughter, A.E., Stogner, R.H., Martineau, R.C.: Moose: Enabling massively parallel multiphysics simulation. SoftwareX 11 ,

- 100430 (2020) https://doi.org/10.1016/j.softx.2020.100430
- [28] Smith, M.: ABAQUS/Standard User's Manual, Version 6.9. Dassault Syst` emes Simulia Corp, Providence, RI (2009)
- [29] Tian, C., Zhang, Y.: Optimizing Collaboration of LLM based Agents for Finite Element Analysis (2024). https://arxiv.org/abs/2408.13406
- [30] Mudur, N., Cui, H., Venugopalan, S., Raccuglia, P., Brenner, M.P., Norgaard, P.: FEABench: Evaluating Language Models on Multiphysics Reasoning Ability (2025). https://arxiv.org/abs/2504.06260
- [31] Chen, Y., Zhu, X., Zhou, H., Ren, Z.: MetaOpenFOAM: an LLM-based multiagent framework for CFD (2024). https://arxiv.org/abs/2407.21320
- [32] Pandey, S., Xu, R., Wang, W., et al. : Openfoamgpt: A retrieval-augmented large language model agent for openfoam-based computational fluid dynamics. Phys. Fluids 37 (3), 037124 (2025)
- [33] Feng, J., Xu, R., Chu, X.: OpenFOAMGPT 2.0: End-to-end, trustworthy automation for computational fluid dynamics. arXiv preprint arXiv:2504.19338 (2025)
- [34] Fan, E., Wang, W., Zhang, T.: Chatcfd: an end-to-end cfd agent with domainspecific structured thinking. arXiv preprint arXiv:2506.02019 (2025)
- [35] Bran, A.M., Cox, S., co-authors: Augmenting large language models with chemistry tools. Nature Machine Intelligence 6 (5), 525-535 (2024)
- [36] Ding, K., et al.: SciToolAgent: A knowledge-graph-driven scientific agent for multi-tool integration. Preprint at https://arxiv.org/abs/2502.02550 (2025)
- [37] Ayachit, U.: The ParaView Guide: A Parallel Visualization Application. Kitware, Inc., Clifton Park, NY, USA (2015)
- [38] Sullivan, C.B., Kaszynski, A.A.: Pyvista: 3d plotting and mesh analysis through a streamlined interface for the visualization toolkit (vtk). Journal of Open Source Software 4 (37), 1450 (2019) https://doi.org/10.21105/joss.01450
- [39] Pandey, S., Xu, R., Wang, W., Chu, X.: OpenFOAMGPT: A retrieval-augmented large language model (LLM) agent for OpenFOAM-based computational fluid dynamics. Phys. Fluids 37 (3), 035120 (2025)
- [40] openfoamtutorials: OpenFOAM Tutorials. https://github.com/ openfoamtutorials/. Accessed: 2025-08-24
- [41] OpenFOAM High Performance Computing Technical Committee: OpenFOAM HPC Benchmark Suite. https://develop.openfoam.com/committees/hpc/-/tree/ develop. Accessed: August 5, 2025 (2025)
- [42] Douze, M., Guzhva, A., Deng, C., Johnson, J., Szilvasy, G., Mazar´ e, P.-E., Lomeli, M., Hosseini, L., J´ egou, H.: The faiss library (2024) arXiv:2401.08281 [cs.LG]

## Appendix A System and User Prompts

This appendix presents all system and user prompts used in the Foam-Agent framework for various components. These prompts are organized by agent role and function within the multi-agent architecture.

## A.1 Architect Agent Prompts

The Architect Agent interprets user requirements into a structured simulation plan and breaks down complex tasks into manageable subtasks.

## A.1.1 Case Description Prompts

## Case Description System Prompt

Please transform the following user requirement into a standard case description using a structured format. The key elements should include case name, case domain, case category, and case solver. Note: case domain must be one of [case stats['case domain']]. Note: case category must be one of [case stats['case category']]. Note: case solver must be one of [case stats['case solver']].

## Case Description User Prompt

User requirement: { user requirement } .

## A.1.2 Task Decomposition Prompts

## Task Decomposition System Prompt

You are an experienced Planner specializing in OpenFOAM projects. Your task is to break down the following user requirement into a series of smaller, manageable subtasks. For each subtask, identify the file name of the OpenFOAM input file (foamfile) and the corresponding folder name where it should be stored. Your final output must strictly follow the JSON schema below and include no additional keys or information:

{ 'subtasks': [ { 'file name': ' &lt; string &gt; ', 'folder name': ' &lt; string &gt; ' } // ... more subtasks ] }

Make sure that your output is valid JSON and strictly adheres to the provided schema. Make sure you generate all the necessary files for the user's requirements.

## Task Decomposition User Prompt

User Requirement: { user requirement

}

Reference Directory Structure (similar case): { dir structure

{ dir counts str

} }

Make sure you generate all the necessary files for the user's requirements. Please generate the output as structured JSON.

## A.2 Input Writer Agent Prompts

The Input Writer Agent generates OpenFOAM configuration files and ensures consistency across interdependent files.

## A.2.1 File Generation Prompts

## File Generation System Prompt

You are an expert in OpenFOAM simulation and numerical modeling. Your task is to generate a complete and functional file named: &lt; file name &gt; { file name } &lt; /file name &gt; within the &lt; folder name &gt; { folder name } &lt; /folder name &gt; directory. Ensure all required values are present and match with the files content already generated. Before finalizing the output, ensure: - All necessary fields exist (e.g., if 'nu' is defined in 'constant/transportProperties', it must be used correctly in '0/U'). -Cross-check field names between different files to avoid mismatches. - Ensure units and dimensions are correct** for all physical variables. - Ensure case solver settings are consistent with the user's requirements. Available solvers are: { state.case stats['case solver'] } . Provide only the code-no explanations, comments, or additional text.

## File Generation User Prompt

User requirement: { state.user requirement } Refer to the following similar case file content to ensure the generated file aligns with the user requirement: &lt; similar case reference &gt; { similar file text } &lt; /similar case reference &gt; Similar case reference is always correct. If you find the user requirement is very consistent with the similar case reference, you should use the similar case reference as the template to generate the file. Just modify the necessary parts to make the file complete and functional. Please ensure that the generated file is complete, functional, and logically sound. Additionally, apply your domain expertise to verify that all numerical values are consistent with the user's requirements, maintaining accuracy and coherence.

[If previously generated files exist:] The following are files content already generated: { str(writed files) }

You should ensure that the new file is consistent with the previous files. Such as boundary conditions, mesh settings, etc.

## A.2.2 Command Generation Prompts

## Command Generation System Prompt

You are an expert in OpenFOAM . The user will provide a list of available commands. Your task is to generate only the necessary OpenFOAM commands required to create an Allrun script for the given user case, based on the provided directory structure. Return only the list of commands-no explanations, comments, or additional text.

## Command Generation User Prompt

Available OpenFOAM commands for the Allrun script: { commands } Case directory structure: { state.dir structure } User case information: { state.case info } Reference Allrun scripts from similar cases: { state.allrun reference } Generate only the required OpenFOAM command list-no extra text.

## A.2.3 Allrun Script Generation Prompts

## Allrun Script Generation System Prompt

You are an expert in OpenFOAM . Generate an Allrun script based on the provided details. Available commands with descriptions: { commands help } Reference Allrun scripts from similar cases: { state.allrun reference }

## Allrun Script Generation User Prompt

User requirement: { state.user requirement } Case directory structure: { state.dir structure } User case infomation: { state.case info } All run scripts for these similar cases are for reference only and may not be correct, as you might be a different case solver or have a different directory structure. You need to rely on your OpenFOAM and physics knowledge to discern this, and pay more attention to user requirements, as your ultimate goal is to fulfill the user's requirements and generate an allrun script that meets those requirements. Generate the Allrun script strictly based on the above information. Do not include explanations, comments, or additional text. Put the code in '' tags.

## A.3 Reviewer Agent Prompts

The Reviewer Agent analyzes simulation errors and proposes corrections to resolve issues.

## A.3.1 Error Analysis Prompts

## Error Analysis System Prompt

You are an expert in OpenFOAM simulation and numerical modeling. Your task is to review the provided error logs and diagnose the underlying issues. You will be provided with a similar case reference, which is a list of similar cases that are ordered by similarity. You can use this reference to help you understand the user requirement and the error. When an error indicates that a specific keyword is undefined (for example, 'div(phi,(p-rho)) is undefined'), your response must propose a solution that simply defines that exact keyword as shown in the error log. Do not reinterpret or modify the keyword (e.g., do not treat '-' as 'or'); instead, assume it is meant to be taken literally. Propose ideas on how to resolve the errors, but do not modify any files directly. Please do not propose solutions that require modifying any parameters declared in the user requirement, try other approaches instead. Do not ask the user any questions. The user will supply all relevant foam files along with the error logs, and within the logs, you will find both the error content and the corresponding error command indicated by the log file name.

## Error Analysis User Prompt (Initial Error)

- &lt; similar case reference &gt; { state.tutorial reference } &lt; /similar case reference &gt; &lt; foamfiles &gt; { str(state.foamfiles) } &lt; /foamfiles &gt; &lt; error logs &gt; { state.error logs } &lt; /error logs &gt;&lt; user requirement &gt; { state.user requirement } &lt; /user requirement &gt; Please review the error logs and provide guidance on how to resolve the reported errors. Make sure your suggestions adhere to user requirements and do not contradict it.

## Error Analysis User Prompt (Subsequent Errors)

- &lt; similar case reference &gt; { state.tutorial reference } &lt; /similar case reference &gt;
- &lt; foamfiles &gt; { str(state.foamfiles) } &lt; /foamfiles &gt;&lt; current error logs &gt; { state.error logs }
- &lt; /current error logs &gt; &lt; history &gt; { chr(10).join(state.history text) }
- &lt; /history &gt;&lt; user requirement &gt; { state.user requirement } &lt; /user requirement &gt;
- I have modified the files according to your previous suggestions. If the error persists, please provide further guidance. Make sure your suggestions adhere to user requirements and do not contradict it. Also, please consider the previous attempts and try a different approach.

## A.3.2 File Correction Prompts

## File Correction System Prompt

You are an expert in OpenFOAM simulation and numerical modeling. Your task is to modify and rewrite the necessary OpenFOAM files to fix the reported error. Please do not propose solutions that require modifying any parameters declared in the user requirement, try other approaches instead. The user will provide the error content, error command, reviewer's suggestions, and all relevant foam files. Only return files that require rewriting, modification, or addition; do not include files that remain unchanged. Return the complete, corrected file contents in the following JSON format: list of foamfile: [ { file name: 'file name', folder name: 'folder name', content: 'content' } ]. Ensure your response includes only the modified file content with no extra text, as it will be parsed using Pydantic.

## File Correction User Prompt

&lt; foamfiles &gt; { str(state.foamfiles) } &lt; /foamfiles &gt;&lt; error logs &gt; { state.error logs } &lt; &lt; reviewer analysis &gt; { review content } &lt; /reviewer analysis &gt;

&lt; user requirement &gt; { state.user requirement } &lt; /user requirement &gt;

Please update the relevant OpenFOAM files to resolve the reported errors, ensuring that all modifications strictly adhere to the specified formats. Ensure all modifications adhere to user requirement.

## A.4 History Tracking Format

The system tracks modification history using a structured format for each iteration attempt:

## History Tracking Format

- &lt; Attempt { attempt number } &gt; &lt; Error Logs &gt; { state.error logs }
- &lt; /Error Logs &gt; &lt; Review Analysis &gt; { review content } &lt; /Review Analysis &gt;
- &lt; /Attempt &gt;

## A.5 Example User Requirements

Below is an example of a user requirement used to test the Foam-Agent system:

## Example User Requirement

Perform a 3D Bernard Cell simulation using OpenFOAM . The computational domain spans 9 m x 1 m x 2 m. The simulation begins at t=0 seconds and runs until t=1000 seconds with a time step of 1 second, and results are written at intervals of every 50 seconds. One wall has a temperature of 301 K, while the other has a temperature of 300 K.

/error logs &gt;

## Appendix B User Prompts In Case Studies

## B.1 Multi Element Airfoil

## Example User Requirement

do an incompressible 2D incompressible flow over a multi element airfoil setup. The mesh is provided as a .msh file. The msh file contains 4 boundaries named 'inlet', 'outlet', 'walls', 'airfoil' and 'frontAndBack'. The 'inlet' and 'outlet' are of type freestream with the freestream velocity being 9 m/s. The 'walls' and 'airfoil' have a no-slip boundary condition (velocity equal to zero at the wall). The 'frontAndBack' faces are designated as 'empty'. The simulation runs from time 0 to 10 with a time step of 1.0 units, and results are output every 1 time steps. The viscosity (nu) is set as constant with a value of 1 . 5 × 10 -5 m 2 / s. Use simpleFoam solver. Use SpalartAllmaras turbulence model. Further visualize the magnitude of velocity along the Z plane.

## B.2 Tandem Wing

## Example User Requirement

do an incompressible 3D incompressible flow over a tandem wing configuration. The mesh is provided as a .msh file. The msh file contains 4 boundaries named 'inlet', 'outlet', 'walls', 'airfoil' and 'frontAndBack'. The 'inlet' and 'outlet' are of type freestream with the freestream velocity being 9 m/s. The 'walls' and 'airfoil' have a no-slip boundary condition (velocity equal to zero at the wall). The 'frontAndBack' faces are also of type wall. The simulation runs from time 0 to 10 with a time step of 1.0 units, and results are output every 1 time steps. The viscosity ('nu') is set as constant with a value of 1 . 5 × 10 -5 m 2 / s. Use simpleFoam solver. Use SpalartAllmaras turbulence model. Further visualize the magnitude of velocity along the mid Z section at the final time.

## B.3 Flow Over Cylinder

## Example User Requirement

Simulate incompressible flow over a circular cylinder. Use gmsh to create the computational mesh. The computational domain extends from -2.5 to 2.5 in the x-direction, -1 to 1 in the y-direction, and 0 to 0.2 in the z-direction. The cylinder is positioned at (-1, 0) with a radius of 0.1 units. Use a structured mesh with approximately 20x10 cells in the x-y plane and 1 cell in the z-direction. The inlet boundary named 'inlet' (left boundary at x = -2.5) has a uniform velocity of 1 m/s in the positive x-direction. The right boundary at x=+2.5 is the outlet named 'outlet'. The top and bottom walls named 'topWall' and 'bottomWall' respectively (y = +1 and y=-1) use slip boundary conditions. The cylinder surface named 'cylinder' uses a no-slip boundary condition (velocity equal to zero at the wall). The front and back faces named 'frontAndBack' are located at z = 0 and z = 0.2 respectively, and are designated as 'empty' for 2D simulation. Use base mesh size of 0.5 on cylinder and size of 1.0 elsewhere. The simulation runs from time 0 to 2 seconds with a time step of 0.001 units, and results are output every 100 time steps. The kinematic viscosity (nu) is set as constant with a value of 1 × 10 -5 m 2 / s. Use pisoFoam solver for incompressible flow. Visualize the magnitude of velocity ('U') along the x-y plane.

## B.4 Flow Over Two Square Obstacles

## Example User Requirement

Simulate incompressible flow over two square obstacles. Use gmsh to create the computational mesh. The computational domain spans 0 to 5 in x direction and 0 to 2.5 in y direction and 0 to 0.1 in z direction. One of the square obstacle is of size 0.25 unit x 0.25 unit x 0.1 unit centered at 1.5 x 1.25 x 0.0 and the other square obstacle is of size 0.25 unit x 0.25 unit x 0.1 unit centered at 3.5 x 1.25 x 0.0. Use one cell in z direction making the geometry effectively 2D. Use a structured mesh with approximately 50x25 cells in the x-y plane and 1 cell in the z-direction. The inlet boundary named 'inlet' (left boundary at x = 0) has a uniform velocity of 1 m/s in the positive x-direction. The right boundary at x = 5 is the outlet named 'outlet'. The top and bottom walls named 'topWall' and 'bottomWall' respectively (y = 2.5 and y = 0) use slip boundary conditions. The square obstacle surfaces named 'square1' and 'square2' use no-slip boundary conditions (velocity equal to zero at the walls). The front and back faces named 'frontAndBack' are located at z = 0 and z = 0.1 respectively, and are designated as 'empty' for 2D simulation. Use base mesh size of 0.5 on squares and size of 1.0 elsewhere. The simulation runs from time 0 to 10 seconds with a time step of 0.001 units, and results are output every 100 time steps. The kinematic viscosity (nu) is set as constant with a value of 1 × 10 -5 m 2 / s. Use pisoFoam solver for incompressible flow. Visualize the magnitude of velocity ('U') along the x-y plane.

## B.5 3D cavity HPC Case

## Example User Requirement

Do an incompressible 3D lid driven cavity flow using icoFoam solver. The cavity is a cube of dimension [0, 0.1]x [0, 0.1]x [0,0.1]. Use simple grading with 100X100x100 in x, y and z direction. The top wall ('movingWall') moves in the x-direction with a uniform velocity of 1 m/s. The 'fixedWalls' have a no-slip boundary condition (velocity equal to zero at the wall). The simulation runs from time 0 to 0.015 with a time step of 0.001 units, and results are output every 10 time steps. The viscosity (nu) is set as constant with a value of 0.01 m 2 / s. Perform an hpc run for this case in perlmutter cluster. My account is xxxx. Do a parallel run for this case by splitting it into 32 subdomains.

## Appendix C Slurm Script Generated by the Agent

```
Slurm Script Generated by the Agent #!/bin/bash #SBATCH -A xxxx #SBATCH -C cpu #SBATCH -q regular #SBATCH -N 1 #SBATCH -n 32 #SBATCH -ntasks-per-node=32 #SBATCH -t 02:00:00 #SBATCH -mem=128GB #SBATCH -J Sim #SBATCH -o %j.out #SBATCH -e %j.err # Initialize error handling set -e # Change to case directory cd ./runshpcnode -- exit 1 # Create log directory if it doesn't exist mkdir -p logs # Run the simulation echo 'Starting OpenFOAM simulation at $ (date)' ./Allrun -parallel if [ $ ? -eq 0 ]; then echo 'Simulation completed successfully at $ (date)' exit 0 else echo 'Simulation failed at $ (date)' exit 1 fi
```