/*
Imposter Word Game (Go + OpenAI)

This is a console-based social deduction game where:
- There are 4 players (1 human, 3 AI).
- One player is randomly assigned as the imposter and does not know the secret word.
- Each player (including human) gives a 1-word clue in a randomized turn order.
- AI voters guess who the imposter is based on all clues given so far.
- The human also votes for who they think is the imposter.
- Majority vote determines whether the imposter is caught (if a tie, randomly chosen).
- Secret words are drawn from a pool and will not repeat until all have been used.
- AI clues are abstract/metaphorical for innocents, vague for imposter.

Features:
- Multiple rounds possible, secret words will not repeat until all have been used.
- Randomized clue and vote order to simulate a real game.
*/

package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// Player represents a single player in the game
type Player struct {
	ID      int
	Name    string
	IsHuman bool
	Role    string // "innocent" or "imposter"
	Clue    string // last clue given by this player
}

var client openai.Client

func main() {
	// Initialize OpenAI client
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("OPENAI_API_KEY not set")
	}
	client = openai.NewClient(option.WithAPIKey(apiKey))

	rand.Seed(time.Now().UnixNano()) // Seed random generator
	reader := bufio.NewReader(os.Stdin)

	// Pool of secret words (ChatGPT generated this list)
	secretWords := []string{
		"Apple", "Mountain", "River", "Galaxy", "Dragon", "Book",
		"Computer", "Chocolate", "Castle", "Ocean", "Music", "Planet",
		"Forest", "Volcano", "Bridge", "Rocket", "Diamond", "Coffee",
		"Mirror", "Sun", "Moon", "Star", "Island",
		"Puzzle", "Tiger", "Camera", "Library", "Garden", "Desert",
	}
	usedWords := map[string]bool{} // Tracks words already used in previous rounds

	playAgain := "y"
	for strings.ToLower(playAgain) == "y" {
		// Build list of words not yet used in this session
		availableWords := []string{}
		for _, w := range secretWords {
			if !usedWords[w] {
				availableWords = append(availableWords, w)
			}
		}

		// Reset used words if we've gone through all of them
		if len(availableWords) == 0 {
			for _, w := range secretWords {
				usedWords[w] = false
			}
			availableWords = append([]string{}, secretWords...)
		}

		// Choose a random word from the available ones
		secretWord := availableWords[rand.Intn(len(availableWords))]
		usedWords[secretWord] = true

		playRound(secretWord, reader) // Play a single round

		fmt.Print("\nDo you want to play again? (y/n): ")
		playAgain, _ = reader.ReadString('\n')
		playAgain = strings.TrimSpace(playAgain)
	}
}

/*
playRound runs a single round of the Imposter Word Game.

- Assigns a random human player and a random imposter.
- Randomizes turn order for clues.
- Records and displays clues in order they were given.
- Conducts voting phase in randomized order.
- Determines whether imposter was caught and prints results.
*/
func playRound(secretWord string, reader *bufio.Reader) {
	playerCount := 4

	// Randomly assign imposter and human player
	imposterIndex := rand.Intn(playerCount)
	humanIndex := rand.Intn(playerCount)

	// Initialize all players
	players := make([]Player, playerCount)
	for i := 0; i < playerCount; i++ {
		players[i] = Player{
			ID:      i,
			Name:    fmt.Sprintf("Player %d", i+1),
			IsHuman: i == humanIndex,
			Role:    "innocent",
		}
	}
	players[imposterIndex].Role = "imposter"

	// Inform human of their role
	if humanIndex == imposterIndex {
		fmt.Println("\nYou are the IMPOSTER. Try to blend in!")
	} else {
		fmt.Println("\nYou are INNOCENT. Secret word is:", secretWord)
	}

	// Clues are stored in the order they were given to display later
	type ClueEntry struct {
		Name string
		Clue string
	}
	var cluesInOrder []ClueEntry
	previousClues := []string{} // Keep track of clues to prevent AI repeats

	// Random turn order for clue giving
	turnOrder := rand.Perm(playerCount)

	// CLUE PHASE
	fmt.Println("\n--- Clue Phase ---")
	for _, i := range turnOrder {
		player := &players[i]

		var clue string
		if player.IsHuman {
			for {
				fmt.Print("Your turn! Enter your 1-word clue: ")
				clueInput, _ := reader.ReadString('\n')
				clue = strings.TrimSpace(clueInput)

				// Validation: non-empty, single word, letters only
				if clue == "" {
					fmt.Println("Clue cannot be empty. Try again.")
					continue
				}
				if strings.Contains(clue, " ") {
					fmt.Println("Please enter only one word. Try again.")
					continue
				}
				valid := true
				for _, r := range clue {
					if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
						valid = false
						break
					}
				}
				if !valid {
					fmt.Println("Clue must contain only letters. Try again.")
					continue
				}

				fmt.Printf("You gave clue: %s\n", clue)
				break
			}
		} else {
			clue = getClueForAI(player.Role, secretWord, i, previousClues)
			fmt.Printf("%s gives clue: %s\n", player.Name, clue)
		}

		player.Clue = clue
		previousClues = append(previousClues, clue)

		// Record clue in order, label human player
		name := player.Name
		if player.IsHuman {
			name += " (YOU)"
		}
		cluesInOrder = append(cluesInOrder, ClueEntry{Name: name, Clue: clue})

		time.Sleep(500 * time.Millisecond) // Small delay for realism
	}

	// Display clues in order of turn
	fmt.Println("\nClues:")
	for _, entry := range cluesInOrder {
		fmt.Printf("%s: %s\n", entry.Name, entry.Clue)
	}

	// VOTING PHASE
	fmt.Println("\n--- Voting Phase ---")
	voteOrder := turnOrder
	votes := []string{}

	for _, i := range voteOrder {
		player := &players[i]

		if player.IsHuman {
			for {
				fmt.Print("Enter who you think is the imposter (Player 1-4): ")
				humanVote, _ := reader.ReadString('\n')
				humanVote = strings.TrimSpace(humanVote)

				// Validate input
				if humanVote == "1" || humanVote == "2" || humanVote == "3" || humanVote == "4" {
					votes = append(votes, humanVote)
					fmt.Printf("You voted for: %s\n", humanVote)
					break
				} else {
					fmt.Println("Invalid input. Please enter a number between 1 and 4.")
				}
			}
		} else {
			vote := getAIVote(players, i, secretWord)
			votes = append(votes, vote)
			fmt.Printf("AI voter %d guesses: %s\n", i+1, vote)
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Count votes
	voteCount := map[string]int{}
	for _, v := range votes {
		voteCount[v]++
	}

	// Determine majority
	majorityVote := ""
	maxVotes := 0
	tiedPlayers := []string{}

	for player, count := range voteCount {
		if count > maxVotes {
			maxVotes = count
			majorityVote = player
			tiedPlayers = []string{player} // new leader, reset tie list
		} else if count == maxVotes {
			tiedPlayers = append(tiedPlayers, player) // another player tied
		}
	}

	// If multiple players are tied, choose one at random
	if len(tiedPlayers) > 1 {
		majorityVote = tiedPlayers[rand.Intn(len(tiedPlayers))]
	}

	realImposter := players[imposterIndex].Name
	fmt.Println("\nMajority vote caught:", majorityVote)

	if majorityVote == realImposter {
		fmt.Println("The imposter was caught!")
		if humanIndex == imposterIndex {
			fmt.Println("You were the imposter. You lose!")
			// Reveal the secret word to the human imposter
			fmt.Println("The secret word was:", secretWord)
		} else {
			fmt.Println("You were innocent. You win!")
		}
	} else {
		fmt.Println("The imposter was NOT caught.")
		fmt.Println("The imposter was:", realImposter)
		if humanIndex == imposterIndex {
			fmt.Println("You were the imposter. You win!")
			// Reveal the secret word to the human imposter
			fmt.Println("The secret word was:", secretWord)
		} else {
			fmt.Println("You were innocent. You lose!")
		}
	}
}

/*
getClueForAI generates a 1-word clue for an AI player.

- Innocent AI: gives abstract, metaphorical, or category-related clues.
- Imposter AI: gives vague, generic clues to blend in.
- Prevents repetition of previous clues in the round.
- Uses GPT API for dynamic, creative clues.
*/
func getClueForAI(role string, secretWord string, aiIndex int, previousClues []string) string {
	var prompt string

	if role == "innocent" {
		personalities := []string{
			"Think poetically and creatively.",
			"Think like a scientist exploring categories.",
			"Think vaguely and abstractly.",
			"Think metaphorically, be unusual.",
		}

		prev := ""
		if len(previousClues) > 0 {
			prev = " Do NOT repeat these clues: " + strings.Join(previousClues, ", ")
		}

		prompt = fmt.Sprintf(
			"You are playing an Imposter game. The secret word is: %s. "+
				"Give ONE subtle, creative, one-word clue that hints at it indirectly. "+
				"Use unusual, abstract, or metaphorical words. %s%s No punctuation. No explanation.",
			secretWord, personalities[aiIndex%len(personalities)], prev)
	} else {
		prev := ""
		if len(previousClues) > 0 {
			prev = " Previous clues: " + strings.Join(previousClues, ", ") +
				". Try to give a plausible word that could blend in."
		}
		prompt = "You are the imposter in a word guessing game. You do NOT know the secret word. " +
			"Give ONE generic, vague word that could plausibly fit multiple topics." + prev + " No punctuation. No explanation."
	}

	resp, err := client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model:       "gpt-4o-mini",
			Temperature: openai.Float(1.0),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
			},
		},
	)

	if err != nil {
		fmt.Println("API Error:", err)
		return "???" // fallback
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content)
}

/*
getAIVote makes an AI player guess who the imposter is.

- AI reads all clues in order.
- Uses GPT to make a plausible guess.
- Returns exactly one player string.
*/
func getAIVote(players []Player, aiIndex int, secretWord string) string {
	prompt := fmt.Sprintf(
		"You are AI voter %d in a social deduction game. One player is the imposter. Think independently.\nHere are the clues:\n", aiIndex+1)

	for _, p := range players {
		name := p.Name
		if p.IsHuman {
			name += " (YOU)"
		}
		prompt += fmt.Sprintf("%s: %s\n", name, p.Clue)
	}

	if players[aiIndex].Role == "innocent" {
		prompt += fmt.Sprintf(
			"\nThe secret word is: %s. Compare each clue to the secret word and vote for the player whose clue seems least related to the secret word.\n", secretWord)
	} else {
		prompt += "\nYou do NOT know the secret word. Vote based on plausibility and weirdness of clues.\n"
	}

	// FORCE THE AI TO RESPOND WITH ONLY ONE PLAYER
	prompt += "IMPORTANT: Respond with exactly ONE player name only: Player 1, Player 2, Player 3, or Player 4. Do NOT include any explanations, reasoning, or extra text."

	resp, err := client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model:       "gpt-4o-mini",
			Temperature: openai.Float(1.0),
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(prompt),
			},
		},
	)

	if err != nil {
		fmt.Println("API Error:", err)
		return "Player ?" // fallback
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content)
}
