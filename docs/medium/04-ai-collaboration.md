# Building JEnv with AI: A Zero-to-Hero Journey

## The Ambitious Non-Go-Developer

The story begins with an idea in 2024.

I wanted to build jenv, a Java version manager specifically optimized for Windows—10 times faster than the existing jenv-for-windows.

But I had a major problem: **I didn't know Go.**

### Traditional Thinking: Learn First, Build Later

Following traditional development logic, I should:

1. Buy "Go Programming Basics" and study from scratch
2. Learn basic syntax (1-2 months)
3. Learn concurrency, file operations, cross-platform development (2-3 months)
4. Only then start writing jenv

**Total time: 3-5 months**—and no guarantee of high-quality code.

### AI Era: Learn While Building

But I'm living in 2024, the **era of AI-assisted development**.

I decided on a different approach:
- **Don't wait for perfection**—start now
- **Learn while building**—AI fills knowledge gaps
- **Iterate rapidly**—from v0.1 to v0.6.7

**Dream big, start small, do it now.**

Result: **6 weeks from zero to releasing v0.6.7.**

This is my AI collaboration journey.

---

## AI Collaboration Chronicle

Over these 6 weeks, I collaborated deeply with 5 AI tools, each playing a critical role at different stages.

Let me walk you through this journey chronologically.

---

### v0.1: Cursor (Rapid Prototyping)

**Timeline**: Weeks 1-2
**Task**: Build basic features (add/list/use/remove)
**Goal**: Quickly validate the idea, create usable v0.1.0

#### Collaboration Style

**My Role**: Product Manager
- I didn't know Go, but I knew what I wanted
- Referenced nvm-windows design, described requirements

**Cursor's Role**: Software Engineer
- Generated Go code framework
- Implemented basic CLI commands
- Handled file I/O and JSON configuration

#### Typical Dialogue

**Me**:
> "I want to implement a `jenv add` command that takes two parameters: alias and JDK path. Then save this mapping to a JSON config file."

**Cursor**:
```go
func AddJDK(alias string, path string) error {
    config := LoadConfig()
    config.JDKs[alias] = path
    return SaveConfig(config)
}
```

**I tested**:
```bash
jenv add java8 "C:\Program Files\Java\jdk1.8.0_291"
jenv list
# Output: java8 -> C:\Program Files\Java\jdk1.8.0_291
```

**I gave feedback**:
> "Great! But if the path doesn't exist, it should show an error."

**Cursor fixed**:
```go
func AddJDK(alias string, path string) error {
    if !fileExists(path) {
        return fmt.Errorf("JDK path not found: %s", path)
    }
    // ... rest of logic
}
```

#### Results

- ✅ Usable v0.1.0 with 4 core commands
- ✅ JSON configuration management
- ✅ Windows path handling

#### Lessons Learned

**✅ What Worked**:
1. **Be specific in requirements**: Not "make an add feature," but "takes two parameters, saves to JSON"
2. **Provide reference examples**: Referenced nvm-windows design to give AI direction
3. **Test immediately**: Run tests immediately, provide instant feedback

**❌ Mistakes Made**:
1. **Vague requirements**: Initially said "make a Java management tool"—Cursor's code was completely off track
2. **Blind trust**: Once didn't test before continuing—later found a bug, wasted time

**Key Insight**:
> **Rapid iteration beats perfect design.** Build first, optimize later.

---

### v0.2: OpenAI GPT-4 (Refactoring & Discussion)

**Timeline**: Week 3
**Task**: Code refactoring, architecture optimization
**Goal**: Make code cleaner and more maintainable

#### Collaboration Style

**My Role**: Code Reviewer
- Pasted Cursor-generated code to GPT-4
- Asked "How can I optimize this code?"

**GPT-4's Role**: Architecture Consultant
- Suggested design patterns
- Pointed out code smells
- Provided best practices

#### Typical Dialogue

**Me**:
> My `config.go` is now 300 lines with global variables everywhere. How do I optimize it?

**GPT-4**:
> Recommend using the **Singleton pattern** for configuration:
> 1. Create `ConfigManager` struct
> 2. Use `sync.Once` to ensure single initialization
> 3. Provide `GetInstance()` method for singleton access
>
> This achieves:
> - Avoids global variables
> - Thread-safe
> - Lazy initialization

**Me**: Got it! Back to Cursor to implement:

```go
type ConfigManager struct {
    config *Config
    mu     sync.RWMutex
}

var (
    instance *ConfigManager
    once     sync.Once
)

func GetInstance() *ConfigManager {
    once.Do(func() {
        instance = &ConfigManager{
            config: loadConfig(),
        }
    })
    return instance
}
```

#### Results

- ✅ Cleaner, modular code structure
- ✅ Applied design patterns (Singleton, Factory)
- ✅ Thread-safe configuration management

#### Lessons Learned

**✅ What Worked**:
1. **Share complete code**: Let GPT-4 see the full picture for better advice
2. **Ask open-ended questions**: "How to optimize" is more valuable than "Is this correct"
3. **Understand before implementing**: Don't copy GPT-4's code directly—understand, then implement with Cursor

**❌ Mistakes Made**:
1. **Over-engineering**: GPT-4 suggested 5 design patterns; I used all of them—code became overly complex
2. **Detached from reality**: Some suggestions were academic but unsuitable for this project's scale

**Key Insight**:
> **AI is a great code reviewer, but you must stay in control.** Not all suggestions should be adopted—choose what fits your project.

---

### v0.3: Trae (UI/UX Enhancement)

**Timeline**: Week 4
**Task**: TUI interface, theme system
**Goal**: Modernize CLI, make it user-friendly

#### Collaboration Style

**My Role**: UX Designer
- Proposed interface requirements: "Want colorful output, dark theme"
- Provided visual references: screenshots of excellent CLI tools

**Trae's Role**: Frontend Engineer
- Integrated Lipgloss library
- Designed color schemes
- Implemented theme switching

#### Typical Dialogue

**Me**:
> I want `jenv list` output to look better—add colors and icons.

**Trae**:
> You can use Charm's Lipgloss library:
>
> ```go
> import "github.com/charmbracelet/lipgloss"
>
> titleStyle := lipgloss.NewStyle().
>     Bold(true).
>     Foreground(lipgloss.Color("#00FF00"))
>
> fmt.Println(titleStyle.Render("📋 Installed Java versions:"))
> ```

**After testing**:
> The green is too bright—can you make it softer?

**Trae adjusted**:
> Changed to `#7CFC00` and added theme system supporting light/dark modes.

#### Results

- ✅ Modern CLI interface
- ✅ Light/dark theme support
- ✅ Colorful output, icons, progress bars

#### Lessons Learned

**✅ What Worked**:
1. **Provide visual references**: Screenshots clearer than text descriptions
2. **Iterate adjustments**: First version doesn't have to be perfect—tweak gradually
3. **UX matters**: Beautiful tools are more appealing to use

**Key Insight**:
> **Don't neglect UI/UX, even for CLI tools.** Good user experience boosts tool appeal.

---

### v0.4: Gemini (Performance Optimization)

**Timeline**: Weeks 5-6
**Task**: Optimize scan performance from 3s to 300ms
**Goal**: Achieve 10x performance improvement

#### Collaboration Style

**My Role**: Performance Engineer
- Shared `scan` function code
- Explained performance issue: "Scanning C: is too slow, takes 3 seconds"

**Gemini's Role**: Optimization Expert
- Analyzed performance bottlenecks
- Suggested concurrency model
- Generated optimized code

#### Typical Dialogue

**Me**:
> My scan function is too slow. How do I optimize it?
>
> ```go
> func Scan(root string) []JDK {
>     var results []JDK
>     filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
>         if strings.HasSuffix(path, "javac.exe") {
>             results = append(results, parseJDK(path))
>         }
>         return nil
>     })
>     return results
> }
> ```

**Gemini**:
> The problem is `filepath.Walk` uses **sequential traversal**—like 1 person searching an entire library.
>
> Recommend using **Dispatcher-Worker concurrency model**:
> 1. Dispatcher divides C: into multiple top-level directories
> 2. Launch multiple Worker goroutines to scan each directory concurrently
> 3. Use channels to collect results
>
> Expected performance improvement: **5-10x**

**Me**: Sounds great! Can you provide a code framework?

**Gemini**:
```go
func ScanConcurrent(root string) []JDK {
    results := make(chan JDK, 10)
    var wg sync.WaitGroup

    dirs := getTopLevelDirs(root)
    for _, dir := range dirs {
        wg.Add(1)
        go func(d string) {
            defer wg.Done()
            scanDir(d, results)
        }(dir)
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    return collectResults(results)
}
```

**I understood the concept**, then went back to Cursor to refine details:
- Added smart filtering (skip `Windows`, `AppData`, etc.)
- Added progress bar display
- Handled edge cases (empty directories, insufficient permissions)

#### Results

- ✅ Scan time: 3s → **300ms** (90% improvement!)
- ✅ Implemented Dispatcher-Worker model
- ✅ Smart filtering, skip useless directories

#### Lessons Learned

**✅ What Worked**:
1. **Quantify the problem**: "Too slow" is less specific than "takes 3 seconds"
2. **Understand concept before implementing**: Gemini explained Dispatcher-Worker model; I understood, then coded
3. **Combine tools**: Gemini provided ideas, Cursor refined details

**❌ Mistakes Made**:
1. **Copy-pasting directly**: First time copied Gemini's code directly—failed on Windows (path issues)
2. **Ignoring edge cases**: Didn't consider insufficient permissions—program crashed

**Key Insight**:
> **AI excels at optimization suggestions, but you control the details.** Understanding principles beats copying code.

---

### All Along: Claude (Holistic Collaboration)

**Timeline**: Throughout development
**Task**: Architecture design, documentation writing, this article
**Goal**: Deep thinking, weighing options

#### Collaboration Style

**My Role**: Product Owner
- Deep conversations with Claude
- Discussing design decisions
- Weighing different approaches

**Claude's Role**: Architecture Consultant + Documentation Assistant
- Provided architecture advice
- Helped write README
- Wrote this article

#### Typical Dialogue

**Me**:
> Should jenv use symbolic links or environment variables?

**Claude**:
> Let's weigh the options:
>
> **Symbolic Links**:
> - ✅ Performance: Zero overhead, direct pointing
> - ✅ Cross-platform: Windows/Linux/macOS all support
> - ❌ Complexity: Needs admin privileges (Windows)
>
> **Environment Variables**:
> - ✅ Simple: No special privileges needed
> - ❌ Performance: Each switch modifies registry/config files
> - ❌ Effectiveness: Needs terminal restart
>
> **Recommendation**: Use symbolic links because:
> 1. Performance is more important (core selling point)
> 2. One-time admin privilege is acceptable
> 3. More cross-platform consistency

**Me**: Makes sense! How to handle permissions?

**Claude**:
> Design it this way:
> 1. Request permissions once during `jenv init`
> 2. All subsequent operations privilege-free
> 3. Give users clear prompts explaining why permissions are needed

**Me**: Perfect!

#### Results

- ✅ Clear architecture decisions
- ✅ Comprehensive README documentation
- ✅ This article (you're reading)

#### Lessons Learned

**✅ What Worked**:
1. **Deep conversations**: Not simple Q&A, but discussion and weighing
2. **Ask "why" often**: Understand reasoning behind decisions
3. **Documentation matters**: Good docs help more people use your tool

**Key Insight**:
> **AI is a thinking partner.** Not just for writing code—more importantly, for helping you think.

---

## AI Collaboration Methodology

After 6 weeks of practice, I've developed an **AI-assisted development methodology**.

### Three-Step Learning for Beginners

If you also want to use AI to learn a new technology, here's my experience:

#### Step 1: AI Generates Code → I Run Tests

- AI generates code framework
- I run it, see results
- If it fails, provide error messages to AI

**Example**:
- Cursor generated `jenv add` function
- I tested, found errors with spaces in paths
- Fed back to Cursor, it added quote handling

#### Step 2: Hit Problems → AI Explains Principles

- Encounter concepts I don't understand during execution
- Ask AI: "What does this mean?"
- AI explains, I understand

**Example**:
- Gemini suggested using goroutines
- I asked: "What's the difference between goroutines and threads?"
- Gemini explained: goroutines are lighter, managed by Go's scheduler
- I understood, continued

#### Step 3: Understand Concept → Modify Code Myself

- Understood the principles
- Try modifying code myself
- If stuck, ask AI again

**Example**:
- Understood Dispatcher-Worker model
- I added smart filtering logic myself
- Didn't know how to add progress bar, asked Cursor—it provided example

### Result

**From knowing zero Go to writing high-performance code—only took 6 weeks.**

Key: **Don't wait for perfection—learn while building.**

---

### Five Tips for Efficient Collaboration

Through practice, I've identified 5 techniques for improving AI collaboration efficiency:

#### Tip 1: Be Specific with Requirements (Provide Reference Examples)

❌ **Bad**: "Make a Java management tool"
✅ **Good**: "Make a tool similar to nvm-windows, but for managing Java versions. Should have add/list/use/remove commands."

**Why**: Specific requirements + reference examples help AI understand what you want.

#### Tip 2: Iterate by Module (Don't Aim for Perfection)

❌ **Bad**: "Make a perfect Java management tool with all features"
✅ **Good**:
- Weeks 1-2: Basic features (add/list/use/remove)
- Week 3: Refactoring
- Week 4: UI enhancement
- Weeks 5-6: Performance optimization

**Why**: Modular iteration lets you see results at each step, adjust direction promptly.

#### Tip 3: Ask Proactively (Don't Pretend to Understand)

❌ **Bad**: See unfamiliar concept (goroutines), pretend to understand, continue
✅ **Good**: Ask immediately: "What are goroutines? How do they differ from threads?"

**Why**: AI won't mind you asking questions. Proactive questioning leads to real learning.

#### Tip 4: Leverage Multiple AIs (Each Has Strengths)

- **Cursor**: Fast coding, great for rapid prototyping
- **GPT-4**: Architecture design, discussing approaches
- **Gemini**: Performance optimization, technical deep dives
- **Trae**: UI/UX, interface design
- **Claude**: Deep thinking, documentation and articles

**Why**: Each AI has its strengths—combining them yields best results.

#### Tip 5: Stay in Control (You're the PM)

❌ **Bad**: Do whatever AI says, completely dependent
✅ **Good**:
- AI provides suggestions
- I weigh pros/cons
- I make final decisions

**Why**: AI is an assistant, not the boss. You're the project owner.

---

### Pitfalls to Avoid

Sharing 4 pitfalls I encountered—hopefully you can avoid them:

#### Pitfall 1: Blindly Trusting AI Code (Always Test)

**My Lesson**:
- Cursor generated `scan` function
- I didn't test, continued to next step
- Later found a bug, wasted half a day

**Correct Approach**:
- After AI generates code, **test immediately**
- Only continue if tests pass
- Test edge cases too (empty input, special characters)

#### Pitfall 2: Copy Without Understanding (You'll Get Stuck)

**My Lesson**:
- Gemini provided concurrent code
- I didn't understand `sync.WaitGroup`, just copied
- Later needed to modify feature—had no idea where to start

**Correct Approach**:
- When you don't understand, **ask first**
- Understand principles, then write code
- Code you don't understand is technical debt

#### Pitfall 3: Vague Requirements (Won't Get What You Want)

**My Lesson**:
- I said "make a scan feature"
- Cursor's generated code was completely wrong
- Went back and forth 5 times

**Correct Approach**:
- **Be specific**: "Scan C: for all `javac.exe`, return JDK path list"
- **Provide reference**: "Similar to nvm-windows scanning"
- **Give examples**: "For example, finding `C:\Program Files\Java\jdk1.8.0_291\bin\javac.exe` should return `C:\Program Files\Java\jdk1.8.0_291`"

#### Pitfall 4: Over-Reliance (Miss Learning Opportunities)

**My Reflection**:
- AI can generate all code
- But if completely reliant, I learn nothing
- Next time facing similar problem, still won't know

**Correct Approach**:
- Use AI to accelerate learning, not replace it
- Understand key concepts yourself
- Try writing code yourself, ask AI when stuck

---

## Perspectives on the AI Era

Through these 6 weeks of development, I've gained profound insights about the AI era.

### Technology Democratization

**Past**:
- Needed years of programming experience to build a project
- Beginners could only start with simple projects, slowly accumulate
- Many ideas stayed as ideas due to high technical barriers

**Now**:
- Have an idea? AI helps you implement it
- Beginners can build complex projects
- **Technology is no longer the barrier—creativity is key**

**Significance**:
- More people can solve their own problems
- More innovative tools emerge
- Technology truly serves people, not just elites

**My Case**:
- I didn't know Go, but built jenv in 6 weeks
- jenv solved my own pain points
- Now helping other Java beginners too

### Personal Customization

**Philosophy**: Everyone can create personalized tools.

**Past**:
- Found a tool not good? Only option was to endure or switch
- Submit an Issue, wait for maintainer to fix—could be months

**Now**:
- Don't like it? Modify it yourself
- AI helps you implement custom features
- **Changing what you don't like has become remarkably easy**

**My Case**:
- jenv-for-windows was too slow, I was unsatisfied
- I didn't wait—built one 10x faster myself
- Now Windows users have a better option

**Encouragement**:
- Don't wait for others to solve your problems
- In the AI era, you can do it yourself
- **Be the change you want to see**

### Dream Big, Start Small

**Dream Big** (Dare to Imagine):
- I wanted to build a tool 10x faster than existing ones
- In the past, this was crazy (don't know Go—how could I?)
- But in the AI era, this is an achievable goal

**Start Small** (Take Small Steps):
- Don't aim for perfection from day one
- Weeks 1-2: Basic features
- Week 3: Refactoring
- Week 4: Enhancement
- Weeks 5-6: Optimization
- v0.1 → v0.6.7, **gradual iteration**

**Do It Now** (Take Immediate Action):
- Don't wait for perfection
- Don't wait to master Go
- **Start now**, learn while building

**Iteration Process**:
- v0.1.0: Basic features, works but slow
- v0.3.0: After refactoring, code cleaner
- v0.4.0: Interface looks better
- v0.6.7: 10x performance boost—release!

**Key Insight**:
> **Done beats perfect.** Release v0.1 first, iterate in use—better than holding back forever.

---

## Closing: You Can Too

By now, you might be thinking:

> "You could do it because you have programming basics. I know zero programming—can I?"

My answer: **Yes!**

### Why You Can Too

1. **AI is more powerful than you think**
   - It can explain every line of code
   - It can teach you from scratch
   - It won't mind how many questions you ask

2. **The key isn't technology—it's ideas**
   - You know your pain points
   - You know what you want
   - AI doesn't know this—only you do

3. **Tools are free/cheap**
   - Cursor: Free version sufficient
   - ChatGPT: Free version very capable
   - Claude/Gemini: Both have free tiers

### Start by Solving Your Own Pain Points

**Don't start with ambitious projects.**

Start with small problems:
- Repetitive daily tasks? Write a script to automate
- Tool not working well? Modify it
- Have a small idea? Give it a try

**jenv started from my pain point:**
- I hated manually editing environment variables
- I was unsatisfied with jenv-for-windows' speed
- I wanted a tool specifically optimized for Windows

### Action Recommendations

If you want to try AI-assisted development too, my advice:

1. **Choose a small project**:
   - Solve your own problem
   - Clear scope, not too large
   - Completable in 1-2 weeks

2. **Choose one AI tool to start**:
   - Cursor: Good for writing code
   - ChatGPT: Good for discussing approaches
   - Pick one, get started

3. **Learn while building**:
   - Don't wait for perfection
   - When facing problems, ask AI
   - Build it, then optimize

4. **Share your story**:
   - Publish on GitHub
   - Write an article
   - Help more people

---

## Start Now

🚀 **Try jenv**:
- GitHub: [github.com/WhyWhatHow/jenv](https://github.com/WhyWhatHow/jenv)
- Landing Page: [jenv-win.vercel.app](https://jenv-win.vercel.app)
- Docs: [Full guide](https://github.com/WhyWhatHow/jenv#usage)

⭐ **Give jenv a star**:
- Let more people see this story
- Encourage me to keep optimizing
- Star count is the best support!

💬 **Share your AI collaboration story**:
- Comment with your experience
- What projects have you built with AI?
- What challenges did you face?

---

**Happy coding! In the AI era, dare to dream, dare to do—create tools that are truly yours! ✨**

---

> **Article Series**:
> 1. [Why I Built JEnv: An Umbrella for Java Beginners](01-why-jenv.md)
> 2. [Zero to Hello World in 5 Minutes with JEnv](02-quick-start.md)
> 3. [Why JEnv is 10x Faster: A Technical Deep Dive](03-technical-deep-dive.md)
> 4. [Building JEnv with AI: A Zero-to-Hero Journey](04-ai-collaboration.md) *(You are here)*
