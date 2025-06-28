# pokedexcli

A terminal-based Pokédex written in Go.  
Explore, catch, and inspect Pokémon using data from the [PokeAPI](https://pokeapi.co). Includes an in-memory caching system and a simple REPL interface.

## 📦 Features

- 🧭 `map` / `mapb`: Paginate through Pokémon location areas  
- 🔍 `explore <area>`: Show Pokémon in a location area  
- 🎯 `catch <pokemon>`: Attempt to catch a Pokémon  
- 📜 `inspect <pokemon>`: View stats, types, height/weight of caught Pokémon  
- 🗂️ `list`: View all Pokémon you’ve caught  
- 💾 Built-in caching to speed up repeated API calls  
- 🆘 `help`: Show all available commands  
- 🚪 `exit`: Quit the app

## 🛠️ Installation

<pre>
  ```bash
git clone https://github.com/kavancamp/pokedexcli.git
cd pokedexcli
go mod tidy
go build
./pokedexcli
</pre>

### 🚀 Usage

Once you run the CLI, you can interact with it like this:
<pre>
Pokedex > help
Pokedex > map
Pokedex > explore route-1
Pokedex > catch pikachu
Pokedex > inspect pikachu
Pokedex > list
Pokedex > exit
</pre>

### 📁 Project Structure

<pre>
pokedexcli/
├── main.go                 # REPL, commands, CLI logic
├── user.go                 # User data model
├── internal/
│   └── pokecache/          # Caching system (map + mutex + expiry)
│       ├── cache.go
│       └── cache_test.go
</pre>

#### 🧠 Technical Highlights

- Uses net/http and encoding/json for API communication
- Caching layer removes stale data using time.Ticker
- CLI supports dynamic command registration via map[string]cliCommand
- Each command is a callback function for modularity

🧑‍💻 Author

Made with ❤️ by Keenah
Feel free to fork, improve, and explore!

