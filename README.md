<p align="center">
  <img src="./logo.png" alt="BeMyEyes Logo" width="200">
</p>
BeMyEyes 👀
Empowering the visually impaired with cutting-edge AI technology.
BeMyEyes is a revolutionary backend application that harnesses the power of artificial intelligence to assist individuals with visual impairments. By combining state-of-the-art machine learning models and seamless integration with the Gemini AI assistant, BeMyEyes provides a comprehensive solution for image recognition, text extraction, and audio transcription.
Features 🚀

Image Recognition: Leverage the capabilities of the ImageNet and Tesseract models to unlock detailed descriptions of visual content, enabling users to comprehend their surroundings with ease.
Optical Character Recognition (OCR): Transform printed materials into accessible digital text with advanced OCR techniques, breaking down barriers to information.
Audio Transcription: Seamlessly convert audio files into readable text using Google's Speech-to-Text (STT) API, ensuring no spoken content goes unheard.
Gemini Integration: Harness the power of Gemini, a cutting-edge AI assistant, to provide natural language responses tailored to the user's needs, making the experience truly intuitive and user-friendly.

API Endpoints 🌐
The BeMyEyes backend exposes the following API endpoints, empowering developers to integrate its capabilities into their applications:

v1/api/imagenet/: Unlock detailed image descriptions by leveraging the ImageNet and Tesseract models.
v1/api/ocr/: Extract text from images with ease, making printed materials accessible to all.
v1/api/gemini/vision: Harness the power of Gemini for visual analysis and response generation.
v1/api/gemini/audio: Transcribe audio files into text format, ensuring no spoken word goes unnoticed.

Getting Started 🚀
Ready to embark on a journey towards inclusivity? Follow these simple steps to set up the BeMyEyes backend locally:

Clone the repository: git clone https://github.com/your-repo/bemyeyes.git
Configure environment variables:

GEMINI_KEY: Your Gemini API key.
GOOGLE_APPLICATION_CREDENTIALS: Path to your Google Cloud Platform service account JSON file for STT API access.

Checkout the build job in the repository to see how you can build and run the project locally.

The build job provides detailed instructions on installing dependencies, building the server, and running it on your local machine. By following the steps in the build job, you'll have the BeMyEyes backend up and running in no time!

Deployment 🚀
BeMyEyes is designed to seamlessly integrate with Google Cloud Run, enabling effortless deployment and scalability. To set up the deployment process, configure the following repository secrets:

ENV_FILE: Base64-encoded contents of the .env file.
GCP_PROJECT_ID: Your Google Cloud Platform project ID.
GCP_SA_KEY: Base64-encoded contents of your Google Cloud Platform service account JSON file.
SERVICE_ACCOUNT_TRANSCRIPTION: Google Cloud Platform service account for the STT API.

Additionally, separate jobs are provided for testing the loading of the OCR, ImageNet, and other models, ensuring robust and reliable performance.
Contributing 🤝
We welcome contributions from the community! If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request. Together, we can make BeMyEyes even better and more inclusive.
License ⚖️
This project is licensed under the MIT License, ensuring its accessibility and encouraging further innovation.
Join us in our mission to empower the visually impaired and create a more inclusive world through the power of artificial intelligence. Let's make a difference, one line of code at a time! 🌟
