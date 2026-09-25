#set document(
  title: [Parallel-processing oriented Programming Language Coa],
	author: ("Ken Shibata"),
	keywords: ("Programming Language", "Experts in Information Science"),
  date: datetime(year: 2025, month: 1, day: 22)
)

Authors: Ken Shibata, #link("https://sites.google.com/site/shinyaty/home")[Shinya Takamaeda]

Mainstream programming languages often require special syntax to improve execution speed through parallel execution, and it is not easy to achieve parallel execution. We have developed a parallel-oriented programming language "Coa" with parallel execution as the default. To avoid race conditions, Coa automatically detects variable dependencies and executes programs in parallel, similar to out-of-order execution of CPUs. Therefore, the behavior from the outside is the same as sequential execution, and it is expected to improve execution speed without increasing program complexity.

- #link("https://github.com/nyiyui/go-coa")[GitHub]
- #link("/projects/coa/coa_swopp-2022.pdf")[Paper for SWoPP 2022]
- #link("https://docs.google.com/presentation/d/1Ooln_pojSlW_T7ZXb0aWG5PyNMTQwvMJoWqgsoFvu0U/preview")[Slides for SWoPP 2022]
- #link("https://youtube.com/watch?v=fqGx2pAEzeg")[Presentation Video (in Japanese)]
