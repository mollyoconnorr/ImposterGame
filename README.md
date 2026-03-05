# Imposter Word Game (Go + OpenAI)

A console-based social deduction game built in **Go** using the **OpenAI GPT API** for AI player behavior. This game has 1 human player and 3 AI players. One player is randomly assigned as the **imposter**, who does **not know the secret word**, while everyone else tries to give subtle clues.

---

##  Game Features

- 4 players: 1 human, 3 AI.
- Randomly assign **imposter** and **human** player each round.
- Human and AI give **1-word clues** in random order.
- AI votes on who they think the imposter is:
  - **Innocent AI**: uses knowledge of the secret word to spot suspicious clues.
  - **Imposter AI**: guesses plausibly to blend in.
- Human votes as well.
- **Majority vote determines** if the imposter is caught. If it's a tie one of those two players is randomly chosen.
- Secret words are drawn from a pool and will **not repeat** until all have been used.
- **Multiple rounds** possible.
- AI clues are abstract/metaphorical for innocents and vague for imposters.

---

##  Requirements

- Go 1.20+ installed
- An **OpenAI API Key**  
  (Sign up at [OpenAI](https://platform.openai.com) if you don’t have one)

  Or email the author, mtmollyoconnor@gmail.com for one! 

---

##  Installation

1. **Clone the repository**

```bash
git clone https://github.com/mollyoconnorr/ImposterGame.git
cd ImposterGame
```

2. **Set your OpenAI API Key**

Mac/Linux:

```bash
export OPENAI_API_KEY="your_api_key_here"
```

3. **Install the OpenAI client using Go Modules**
   
```bash
go mod tidy
```

3. **Run the game**

```bash
go run main.go
```

---

##  How to Play

1. The game will inform you if you are **innocent** or the **imposter**.
2. If innocent, you will see the **secret word**.
3. Each player takes a turn giving a **1-word clue**:
   - **Human**: type your clue when prompted.
   - **AI**: will automatically generate a clue.
4. After all clues are given, everyone votes on who they think the **imposter** is.
5. The game announces whether the **imposter was caught**.
6. Play again as many rounds as you like — words **do not repeat** until all have been used.

---

##  Notes

- Human players must enter **single-word clues** using letters only.
- Votes must be a number between **1 and 4** corresponding to the players.
- If the human is the imposter, the **secret word will be revealed** at the end of the round.
- AI voting is powered by GPT, so votes may vary depending on clues.

---

##  Example

```bash
You are INNOCENT. Secret word is: Diamond

--- Clue Phase ---
Player 3 gives clue: Lustre
Player 4 gives clue: Eternity
Your turn! Enter your 1-word clue: ring
You gave clue: ring
Player 1 gives clue: Sparkle

--- Voting Phase ---
AI voter 4 guesses: Player 3
AI voter 1 guesses: Player 1
Enter who you think is the imposter (Player 1-4): 4
You voted for: 4

Majority vote caught: Player 4
The imposter was NOT caught.
The imposter was: Player 1
You were innocent. You lose!
```

---

##  Customization

- Add or remove words from the `secretWords` list in `main.go`.
- Adjust `playerCount` if you want more AI players.
- Modify AI prompt logic in `getClueForAI` or `getAIVote` to tweak AI behavior.

---
