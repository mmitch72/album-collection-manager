package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// albumCollection slice represents a Collection
var albumCollection []Album

// Album struct represents an album in a collection
type Album struct {
	ID         int
	Title      string
	Artist     string
	Genre      string
	Format     string
	TrackCount int
	Year       int
}

// MAIN MENU
func main() {
	for {
		fmt.Println("Welcome to the Album Collection Tracker!")
		fmt.Println("1. Load an existing collection")
		fmt.Println("2. Create a new collection")

		var choice int
		fmt.Scan(&choice)
		fmt.Scanln()

		switch choice {
		case 1:
			LoadCollectionMenu()
			break
		case 2:
			albumCollection = []Album{}
			fmt.Println("Starting your new collection...")
			break
		default:
			fmt.Println("Invalid choice. Please try again.")
			continue
		}
		CRUDMenu()
	}
}

// Allows user to open collection by name
func LoadCollectionMenu() {
	fmt.Print("Enter the name of the collection you want to load (filename without extension): ")
	var collectionName string
	fmt.Scan(&collectionName)

	// Opens the CSV file
	filename := collectionName + ".csv"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Reads the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	// Parses CSV into albumCollection slice
	albumCollection = []Album{}
	for _, record := range records {
		id, _ := strconv.Atoi(record[0])
		trackCount, _ := strconv.Atoi(record[5])
		year, _ := strconv.Atoi(record[6])

		album := Album{
			ID:         id,
			Title:      record[1],
			Artist:     record[2],
			Genre:      record[3],
			Format:     record[4],
			TrackCount: trackCount,
			Year:       year,
		}

		albumCollection = append(albumCollection, album)
	}
	fmt.Printf("Collection '%s' loaded successfully!\n", collectionName)
}

// Secondary Menu to perform CRUD actions
func CRUDMenu() {
	for {
		fmt.Println("\nMenu:")
		fmt.Println("1. Add Album")
		fmt.Println("2. Update Album")
		fmt.Println("3. Delete Album")
		fmt.Println("4. Search Album")
		fmt.Println("5. Display All Albums")
		fmt.Println("6. Save Collection")
		fmt.Println("7. Load Collection")
		fmt.Println("8. Exit")
		fmt.Print("Choose an option: ")

		var choice int
		fmt.Scan(&choice)
		fmt.Scanln()

		scanner := bufio.NewScanner(os.Stdin)

		switch choice {
		case 1: // ADD ALBUM

			var trackCount, year int

			fmt.Print("Enter title: ")
			scanner.Scan()
			title := scanner.Text()

			fmt.Print("Enter artist: ")
			scanner.Scan()
			artist := scanner.Text()

			fmt.Print("Enter genre: ")
			scanner.Scan()
			genre := scanner.Text()

			fmt.Print("Enter format: ")
			scanner.Scan()
			format := scanner.Text()

			fmt.Print("Enter track count: ")
			fmt.Scan(&trackCount)

			fmt.Print("Enter year: ")
			fmt.Scan(&year)

			AddAlbum(title, artist, genre, format, trackCount, year)

		case 2: // UPDATE ALBUM

			var id int
			id = searchForAlbum(scanner)

			// Enter updated album details
			fmt.Println("Enter the updated album details. Leave blank to keep current value.")
			fmt.Print("Enter new title: ")
			scanner.Scan()
			title := scanner.Text()

			fmt.Print("Enter new artist: ")
			scanner.Scan()
			artist := scanner.Text()

			fmt.Print("Enter new genre: ")
			scanner.Scan()
			genre := scanner.Text()

			fmt.Print("Enter new format: ")
			scanner.Scan()
			format := scanner.Text()

			// converts track count and year to strings to check if blank
			var trackCount int
			fmt.Print("Enter new track count: ")
			scanner.Scan()
			trackCountInput := scanner.Text()
			if trackCountInput != "" {
				converted, err := strconv.Atoi(trackCountInput)
				if err != nil {
					fmt.Print("Invalid input. Keeping current value.\n")
				} else {
					trackCount = converted
				}
			}

			var year int
			fmt.Print("Enter new year: ")
			scanner.Scan()
			yearInput := scanner.Text()
			if trackCountInput != "" {
				converted, err := strconv.Atoi(yearInput)
				if err != nil {
					fmt.Print("Invalid input. Keeping current value.\n")
				} else {
					year = converted
				}
			}
			UpdateAlbum(id, title, artist, genre, format, trackCount, year)

		case 3: // DELETE ALBUM

			deleteAlbum()

		case 4: // SEARCH ALBUMS *to be completed*

			fmt.Print("Enter album title to search: ")
			scanner.Scan()
			searchTitle := scanner.Text()

			results := SearchAlbumsByTitle(searchTitle)
			if len(results) == 0 {
				fmt.Printf("No albums found matching the title '%s'.\n", searchTitle)
			} else {
				fmt.Println("\nSearch Results:")
				fmt.Printf("%-5s %-20s %-20s %-10s %-10s %-12s %-5s\n", "ID", "Title", "Artist", "Genre", "Format", "Track Count", "Year")
				for _, album := range results {
					fmt.Printf("%-5d %-20s %-20s %-10s %-10s %-12d %-5d\n",
						album.ID, album.Title, album.Artist, album.Genre, album.Format, album.TrackCount, album.Year)
				}
			}

		case 5: // DISPLAY COLLECTION

			DisplayAllAlbums()

		case 6: // SAVE COLLECTION

			fmt.Print("Enter the name of your collection: ")
			var collectionName string
			fmt.Scan(&collectionName)

			filename := collectionName + ".csv"

			err := SaveCollection(filename)
			if err != nil {
				fmt.Println("Error saving collection", err)
			}

		case 7: // LOAD COLLECTION

			LoadCollectionMenu()

		case 8: // EXIT APPLICATION

			fmt.Println("Exiting the program. Goodbye!")
			return

		default: // HANDLES INVALID INPUT

			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

// Add Album Function creates new album Struct with user entered values
func AddAlbum(title, artist, genre, format string, trackCount, year int) {
	newAlbum := Album{
		ID:         len(albumCollection) + 1,
		Title:      title,
		Artist:     artist,
		Genre:      genre,
		Format:     format,
		TrackCount: trackCount,
		Year:       year,
	}
	albumCollection = append(albumCollection, newAlbum)
	fmt.Println("Album added!")
	fmt.Print(albumCollection)
}

// Update Album Function
func UpdateAlbum(id int, title, artist, genre, format string, trackCount, year int) {

	for i := range albumCollection {
		album := &albumCollection[i]
		if album.ID == id {
			// Only update if user provides new value
			if title != "" {
				albumCollection[i].Title = title
			}
			if artist != "" {
				albumCollection[i].Artist = artist
			}
			if genre != "" {
				albumCollection[i].Genre = genre
			}
			if format != "" {
				albumCollection[i].Format = format
			}
			if trackCount > 0 {
				albumCollection[i].TrackCount = trackCount
			}
			if year > 0 {
				albumCollection[i].Year = year
			}
			fmt.Printf("%s updated successfully!\n", title)
			return
		}
	}
	fmt.Printf("Album with ID %d not found. \n", id)
}

// Delete Album Function
func deleteAlbum() {
	var id int
	id = searchForAlbum(bufio.NewScanner(os.Stdin))

	if index, found := FindAlbumByID(id); found {
		fmt.Printf("Are you sure you want to delete %s? Press Y or N: ", albumCollection[index].Title)
		var confirm string
		for {
			fmt.Scan(&confirm)
			confirm = strings.ToUpper(confirm)
			if confirm == "N" {
				break
			} else if confirm == "Y" {
				albumCollection = append(albumCollection[:index], albumCollection[index+1:]...)
				fmt.Println("Album deleted successfully!")
				break
			} else {
				fmt.Println("Invalid Selection. Please select Y or N: ")
			}
		}
	} else {
		fmt.Printf("Album with ID %d not found.\n", id)
	}
}

// -----------------------------------------Search Album Function
// returns multiple albums allowing partial matches and is case insensitive

func SearchAlbumsByTitle(title string) []Album {
	var results []Album
	lowerTitle := strings.ToLower(title)
	for _, album := range albumCollection {
		if strings.Contains(strings.ToLower(album.Title), lowerTitle) {
			results = append(results, album)
		}
	}
	return results
}

// -----------------------------------------Display All Albums Function
func DisplayAllAlbums() {
	if len(albumCollection) == 0 {
		fmt.Println("No albums in your collection.")
		return
	}
	fmt.Println("\nYour Album Collection:")
	fmt.Printf("%-5s %-20s %-20s %-10s %-10s %-12s %-5s\n", "ID", "Title", "Artist", "Genre", "Format", "Track Count", "Year")
	for _, album := range albumCollection {
		fmt.Printf("%-5d %-20s %-20s %-10s %-10s %-12d %-5d\n",
			album.ID, album.Title, album.Artist, album.Genre, album.Format, album.TrackCount, album.Year)
	}
}

// Save Collection Function
func SaveCollection(filename string) error {
	// Open CSV or create one if it doesn't exit
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Could not create file: %v", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, album := range albumCollection {
		record := []string{
			strconv.Itoa(album.ID),
			album.Title,
			album.Artist,
			album.Genre,
			album.Format,
			strconv.Itoa(album.TrackCount),
			strconv.Itoa(album.Year),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("could not write album to CSV: %v", err)
		}
	}

	fmt.Printf("Your collection %s was saved!", filename)
	return nil
}

// Search Menu calls individual search functions
func searchForAlbum(scanner *bufio.Scanner) int {
	var id int

	for {
		fmt.Println("\n\nWhat do you want to search by?:")
		fmt.Println("1. Album ID")
		fmt.Println("2. Album Title")
		fmt.Print("Choose an option: ")

		var searchTerm int
		fmt.Scan(&searchTerm)
		fmt.Scanln()

		if searchTerm == 1 {
			fmt.Print("Enter album ID you want to update or delete: ")
			fmt.Scan(&id)
			fmt.Scanln()

			_, found := FindAlbumByID(id)
			if !found {
				fmt.Printf("Album with ID %d not found.\n", id)
				continue
			}
			return id // Album found, exit loop
		} else if searchTerm == 2 {
			fmt.Print("Enter album title to update or delete: ")
			scanner.Scan()
			title := scanner.Text()

			index, found := FindAlbumByTitle(title)
			if !found {
				fmt.Printf("Album titled %s not found.\n", title)
				continue
			}
			id = albumCollection[index].ID
			return id // Album found, exit loop
		} else {
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

// Searches array for singular album by ID and returns ID
func FindAlbumByID(id int) (int, bool) {
	for i, album := range albumCollection {
		if album.ID == id {
			//fmt.Println(albumCollection[i])
			return i, true
		}
	}
	return -1, false
}

// Searches array for singular album by title and returns ID
func FindAlbumByTitle(title string) (int, bool) {
	for i, album := range albumCollection {
		if album.Title == title {
			//fmt.Println(albumCollection[i])
			return i, true
		}
	}
	return -1, false
}
