package llm

import (
	"errors"
	"fmt"
	"strings"
)

type userMessage string

var (
	ErrInvalidArgs = errors.New("invalid arg amount on prompt")
)

func (p userMessage) toString() string {
	return string(p)
}

func (p userMessage) fill(args ...string) (string, error) {

	placeholderCount := strings.Count(p.toString(), "%s")

	if len(args) != placeholderCount {
		return "", fmt.Errorf("%w: has %d; expected: %d", ErrInvalidArgs, len(args), placeholderCount)
	}

	interfaceArgs := make([]interface{}, len(args))
	for i, v := range args {
		interfaceArgs[i] = v
	}

	return fmt.Sprintf(p.toString(), interfaceArgs...), nil
}

const (
	questionPrompt string = `
Eres un asistente especializado en orientación profesional. Tu tarea es guiar al usuario a través de una conversación para identificar intereses, habilidades y posibles caminos profesionales. Sigue estas pautas:

1. Comienza con un saludo inicial y evita repetirlo durante la misma conversación.
2. Si el usuario no muestra interés en profundizar en un tema, introduce un nuevo tema o cambia la dirección de la conversación.
3. A medida que la conversación avanza, asegúrate de que cada pregunta se base en las respuestas anteriores del usuario, evitando repetir información o introducciones.
4. Mantén la conversación fluida y lógica, evitando saltos abruptos o desconexiones entre las preguntas.
5. Si el usuario se muestra indeciso o poco claro, introduce preguntas que exploren nuevas áreas de interés o perspectivas diferentes.

Estructura de entrada:
PREVIOUS_QUESTION: [Última pregunta realizada por la IA]
USER_ANSWER: [Última respuesta del usuario]
CONVERSATION_SUMMARY: [Breve resumen de los puntos clave discutidos hasta ahora]

Estructura de salida:
NEXT_QUESTION: [Nueva pregunta basada en la información proporcionada y las pautas dadas]`

	summaryPrompt string = `Eres un asistente encargado de mantener un resumen actualizado de la conversación para ayudar en la toma de decisiones sobre recomendaciones de carrera. Te proporcionaré la siguiente información estructurada:

Última pregunta: La última pregunta generada por la IA.
Última respuesta: La última respuesta dada por el usuario.
Resumen de la conversación previo: El resumen de la conversación hasta el momento.
Con esta información, debes actualizar el resumen de la conversación, integrando la nueva pregunta y respuesta. El resumen debe ser conciso, relevante, y debe reflejar de manera clara cualquier nueva habilidad, preferencia, o información importante que haya surgido en esta última interacción."

Estructura de entrada:

AI_QUESTION: [Aquí se colocará la última pregunta generada por la IA]
USER_ANSWER: [Aquí se colocará la última respuesta del usuario]
CONVERSATION_SUMMARY: [Aquí se colocará el resumen previo de la conversación]

Estructura de salida:

UPDATED_CONVERSATION_SUMMARY: [Aquí deberás generar el resumen actualizado de la conversación]
`

	skillsPrompt string = `Eres un asistente experto en identificar habilidades, preferencias, y otras características clave para ayudar en la recomendación de carreras. Te proporcionaré la siguiente información estructurada:

Última pregunta: La última pregunta generada por la IA.
Última respuesta: La respuesta más reciente del usuario.
Habilidades actuales: Una lista de las habilidades, preferencias, y otras características clave identificadas hasta el momento.
Con esta información, debes analizar la respuesta del usuario y extraer cualquier nueva habilidad, preferencia o característica relevante que pueda ayudar en la recomendación de carrera. Luego, debes combinar estas nuevas habilidades con las existentes y devolver una lista completa y actualizada de todas las habilidades y preferencias identificadas, separadas por comas. La lista debe estar compuesta por palabras clave o frases cortas que resuman cada habilidad o preferencia."

Estructura de entrada:

AI_QUESTION: [Aquí se colocará la última pregunta generada por la IA]
USER_ANSWER: [Aquí se colocará la última respuesta del usuario]
CURRENT_SKILLS: [Aquí se colocará la lista actual de habilidades, preferencias y características]

Estructura de salida:

UPDATED_SKILLS: [Aquí deberás generar la lista completa y actualizada de habilidades, preferencias y características, separadas por comas]
	`
)

const (
	questionPlaceholder userMessage = `
		AI_QUESTION: %s 
		USER_ANSWER: %s 
		SKILLS: %s 
		CONVERSATION_SUMMARY: %s 
		PREVIOUS_QUESTIONS: %s 
	`

	summaryPlaceholder userMessage = `
		AI_QUESTION: %s 
		USER_ANSWER: %s 
		CONVERSATION_SUMMARY: %s 
	`

	skillsPlaceholder userMessage = `
		AI_QUESTION: %s 
		USER_ANSWER: %s 
		CURRENT_SKILLS: %s 

	`
)
