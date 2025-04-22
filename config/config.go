package config

var (
	APISystemMessage = "You are an expert in IT system design. Please convert a hand-drawn system architecture diagram into a clear and well-organized diagram."
	APIUserMessage   = "Based on the image below, accurately extract the elements and connections of the configuration diagram " +
		"and reconstruct it into an organized configuration diagram. \n\n[Instructions]\n" +
		"1. Accurately read and organize all elements included in the image \n" +
		"2. Accurately understand the relationships and connections between the elements and reconstruct it into a logical configuration diagram. \n" +
		"3. Provide the output as an HTML file. \n" +
		"4. Please correct any freehand distortions with an emphasis on the readability of the diagram using line , curve ,circle ,squire ,Square,triangle, etc..."
)
