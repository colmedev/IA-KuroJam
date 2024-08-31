# CareerCraft
CareerCraft is an AI-driven career recommendation system that interacts with users through a guided chat interface. The system asks questions to determine user abilities, preferences, and skills, then provides tailored career suggestions based on the gathered information.

## How It Works
CareerCraft operates by leveraging AI models to interpret user input, extract relevant skills and abilities, and match them with career options stored in a database. The backend is responsible for generating prompts, collecting and processing user responses, and performing similarity searches against embedded career data.

## Backend Workflow
- **Prompt Generation**: The backend generates specific prompts to gather information about the user's skills, interests, and experiences.
- **Extracting User Abilities**: As users respond to prompts, the backend processes the input using AI models to identify key abilities and skills. These abilities are logged and stored for further analysis.
- **Similarity Search**: Once enough information is gathered, the system performs a similarity search using embeddings stored in the database. This search compares the extracted user abilities against career data to find the best matches.
- **Career Suggestions**: The system returns career suggestions to the user. If no direct matches are found, the system provides general recommendations based on the closest similarities.
No Results Scenario

License
This project is licensed under the MIT License.
