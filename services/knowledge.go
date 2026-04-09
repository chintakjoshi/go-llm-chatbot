package services

// Project represents a single project with all its details
type Project struct {
	Name         string
	Description  string
	Backstory    string
	Technologies []string
	Features     []string
	Links        map[string]string
	Category     string
}

// PersonalInfo contains all personal and professional information
type PersonalInfo struct {
	Name         string
	Title        string
	Experience   string
	CurrentRole  string
	Education    []string
	Achievements []string
	Skills       map[string][]string
	Links        map[string]string
}

// GetPersonalInfo returns structured personal information
func GetPersonalInfo() PersonalInfo {
	return PersonalInfo{
		Name:        "Chintak Joshi",
		Title:       "Full-Stack Engineer",
		Experience:  "4+ years of experience building and designing highly scalable, reliable services",
		CurrentRole: "United Rentals",
		Education: []string{
			"Saint Louis University - Master of Science in Computer Science",
			"Gujarat Technological University - Bachelor's of Engineering in Information Technology",
		},
		Achievements: []string{
			"Proud member of Open Source with SLU, where I led open-source projects across various stacks",
			"Oracle Cloud Certified in OCI Developer Professional, Application Integration Professional",
			"Oracle Cloud Certified AI Foundations, Fusion AI Agent Studio Foundations, Data Platform Foundations",
		},
		Skills: map[string][]string{
			"Languages":      {"Python", "JavaScript", "TypeScript", "Java", "Go", "PHP", "HTML", "CSS", "Shell scripting"},
			"Frameworks":     {"React", "Next.js", "Node.js", "Django", "FastAPI", "Flask", "Spring Boot", "Gin"},
			"AI/ML":          {"TensorFlow", "PyTorch", "scikit-learn", "Keras", "Pandas", "OpenCV", "AWS Textract", "MCP", "LangChain", "Weaviate", "Pinecone"},
			"Databases":      {"MySQL", "PostgreSQL", "MongoDB", "DynamoDB", "Vector databases"},
			"DevOps & Cloud": {"AWS (EC2, S3, RDS, Lambda)", "Bazel", "Helm", "Argo CD", "Docker", "Git", "Jenkins", "Terraform", "Kubernetes"},
		},
		Links: map[string]string{
			"LinkedIn":  "https://www.linkedin.com/in/chintakjoshi",
			"GitHub":    "https://github.com/chintakjoshi",
			"Portfolio": "https://chintakjoshi.github.io/chintakjoshi/",
		},
	}
}

// GetProjects returns all actual projects
func GetProjects() []Project {
	return []Project{
		{
			Name:         "iSpraak",
			Description:  "I built iSpraak to help language learners practice speaking and get instant corrective feedback.",
			Backstory:    "This started as a research project at SLU. The idea of using speech tech to actually help people learn pronunciation felt way more useful than a typical class assignment, so I ran with it.",
			Technologies: []string{"PHP", "React", "Amazon Polly (AWS TTS)", "Google OAuth 2.0", "Tailwind CSS", "Docker", "AWS", "S3", "EC2"},
			Features: []string{
				"Automated speech evaluation",
				"Instant corrective feedback",
				"Responsive design with modern UI",
			},
			Links: map[string]string{
				"Live Demo":   "https://ispraak.net/",
				"Source Code": "https://github.com/dnickol1/ispraak_open/",
			},
			Category: "Academic/Personal",
		},
		{
			Name:         "Orbital Sweep",
			Description:  "I built Orbital Sweep as a browser-friendly 3D game where players fly through space and clean up orbital debris.",
			Backstory:    "I wanted to learn Three.js and matrix math, so I figured making a space game would be more fun than reading docs. First time doing anything 3D in a browser and it was a blast.",
			Technologies: []string{"Three.js", "gl-matrix", "JavaScript", "Vite"},
			Features: []string{
				"3D gameplay focused on collecting orbital debris",
				"Browser and mobile-friendly rendering loop",
				"Learning project centered on scene setup, animation, and matrix math",
			},
			Links: map[string]string{
				"Play Game":   "https://rolling-ball-three-js-chintak-joshis-projects.vercel.app/",
				"Source Code": "https://github.com/chintakjoshi/orbital-sweep",
			},
			Category: "Personal",
		},
		{
			Name:         "Token-based Authentication App",
			Description:  "I built this authentication app around token-based login flows and OTP-based email verification.",
			Backstory:    "Auth is one of those things every app needs but nobody wants to build from scratch. I wanted to understand JWT flows deeply instead of just plugging in a library, so I built the whole thing end to end.",
			Technologies: []string{"React", "Spring Boot", "PostgreSQL", "Docker", "Java", "TypeScript"},
			Features: []string{
				"Token-based authentication",
				"Email verification via OTP",
				"Spring Boot backend with React frontend",
			},
			Links: map[string]string{
				"Source Code": "https://github.com/chintakjoshi/authapp",
			},
			Category: "Academic",
		},
		{
			Name:         "Transformative Workforce Academy",
			Description:  "I contributed to this open-source platform to help justice-involved job seekers and second-chance employees find work opportunities.",
			Backstory:    "This was part of Open Source with SLU. Working on something that actually helps people in the St. Louis community find jobs made the late-night coding sessions feel worth it.",
			Technologies: []string{"React", "Node.js", "Express.js", "PostgreSQL", "Docker"},
			Features: []string{
				"Hiring platform for justice-involved talent",
				"Promotes self-sufficiency and public safety",
				"Open source project for St. Louis community",
			},
			Links: map[string]string{
				"Source Code": "https://github.com/chintakjoshi/TWA-OSS",
			},
			Category: "Open Source",
		},
		{
			Name:         "Sync Playlist",
			Description:  "I built Sync Playlist to help users move and sync playlists across services like Spotify and YouTube Music.",
			Backstory:    "I kept losing playlists when switching between Spotify and YouTube Music. Instead of doing it manually every time, I built a tool to handle it. Scratching my own itch, basically.",
			Technologies: []string{"Next.js", "React", "TypeScript", "Go", "Gin", "PostgreSQL", "Docker", "JWT", "Google OAuth2"},
			Features: []string{
				"Cross-platform playlist transfer and sync flows",
				"Intelligent track matching with confidence scoring and transfer history",
				"OAuth2 authentication, token refresh, and rate-limited API integration",
			},
			Links: map[string]string{
				"Source Code": "https://github.com/chintakjoshi/sync-playlist",
			},
			Category: "Personal",
		},
		{
			Name:         "Online Grocery Store for SLU",
			Description:  "I led this SLU grocery store project to manage stock and support online ordering.",
			Backstory:    "This was a team project at SLU and I ended up leading it. Managing a team while shipping features taught me more about coordination than any class could.",
			Technologies: []string{"React", "Node.js", "Express.js", "PostgreSQL", "Docker"},
			Features: []string{
				"Grocery stock management",
				"Online ordering system",
				"Team leadership role",
			},
			Links: map[string]string{
				"Source Code": "https://github.com/chintakjoshi/onlinestore",
			},
			Category: "Academic",
		},
		{
			Name:         "What Is Around Me!",
			Description:  "I built this location-based app to help users discover nearby places on a map and in a searchable list.",
			Backstory:    "I was new to React Native and wanted a project that would force me to deal with maps, GPS, and real API data all at once. Turned out to be a great way to learn mobile development.",
			Technologies: []string{"React Native", "Expo", "TypeScript", "Node.js", "Express", "Google Places API", "react-native-maps"},
			Features: []string{
				"Current-location discovery with nearby search filters",
				"Map and list experiences backed by a shared mobile context",
				"Place detail views with ratings, hours, contact info, and navigation links",
			},
			Links: map[string]string{
				"Source Code": "https://github.com/chintakjoshi/WhatIsAroundMe",
			},
			Category: "Personal",
		},
		{
			Name:         "Collaborative Draw Board",
			Description:  "I built this real-time collaborative drawing board so multiple users can join a board and sketch together live.",
			Backstory:    "I was curious about WebSockets and real-time sync, so I built something where you can actually see the results instantly. Drawing with friends in real time made the debugging way more entertaining.",
			Technologies: []string{"FastAPI", "React", "TypeScript", "WebSockets", "SQLite", "Tailwind CSS", "React Konva"},
			Features: []string{
				"Multi-user drawing with live strokes, shapes, text, and cursors",
				"Board creation and joining through shareable 6-character codes",
				"SQLite persistence for board state, layers, and drawing history",
			},
			Links: map[string]string{
				"Source Code": "https://github.com/chintakjoshi/collaborative-drawing-board",
			},
			Category: "Personal",
		},
		{
			Name:         "E-verification of Documents",
			Description:  "I built this document verification system to validate records and generate QR-based identity checks.",
			Backstory:    "This was my undergrad project. The problem was simple but real: how do you verify a document is legit without calling someone? QR codes turned out to be a clean solution.",
			Technologies: []string{"Django", "MySQL", "ORM"},
			Features: []string{
				"Document validation against system dataset",
				"Unique hex code generation",
				"QR code for identity verification",
			},
			Links: map[string]string{
				"Source Code": "https://github.com/chintakjoshi/QR-project",
			},
			Category: "Academic",
		},
		{
			Name:         "Autonomous Research Assistant for Scientific Literature",
			Description:  "I built this multi-agent research assistant to scan and summarize papers from arXiv and PubMed in real time.",
			Backstory:    "Reading research papers is a time sink. I wanted to see if I could get multiple AI agents to divide and conquer the work, and the multi-agent pattern turned out to be a really interesting problem to solve.",
			Technologies: []string{"React", "Flask", "FastAPI", "MCP", "LangChain", "Pinecone"},
			Features: []string{
				"Real-time paper scanning from arXiv and PubMed",
				"Multi-agent system with planner-executor framework",
				"Task decomposition and role delegation",
			},
			Links: map[string]string{
				"Source Code": "https://github.com/chintakjoshi/auto_research_assis",
			},
			Category: "AI/ML",
		},
		{
			Name:         "RAG with Neural Retrieval",
			Description:  "I built this RAG pipeline to experiment with learned retrieval, reranking, local generation, and evaluation end to end.",
			Backstory:    "Everyone talks about RAG but few people build the retrieval part from scratch. I wanted to go deeper than just plugging in an embedding model and actually understand bi-encoders, rerankers, and evaluation.",
			Technologies: []string{"Python", "FastAPI", "PyTorch", "TensorFlow", "FAISS", "MLflow", "Ollama", "BEIR"},
			Features: []string{
				"Bi-encoder retrieval pipeline with FAISS search and BM25 baselines",
				"Cross-encoder, ColBERT-style, and distilled reranking workflows",
				"FastAPI serving, MLflow experiment tracking, and local grounded generation",
			},
			Links: map[string]string{
				"Source Code": "https://github.com/chintakjoshi/Custom-RAG-Pipeline-with-Neural-Retrieval",
			},
			Category: "AI/ML",
		},
		{
			Name:         "Warfarin Dose Prediction System",
			Description:  "I built this machine learning system to predict Warfarin dosage from 15+ clinical features and share it through a Gradio app.",
			Backstory:    "This was a cool intersection of ML and healthcare. Warfarin dosing is tricky because it varies so much between patients, so predicting it with ML felt like a meaningful problem to tackle.",
			Technologies: []string{"React", "Flask", "scikit-learn", "Keras", "TensorFlow", "Pandas", "NumPy", "Gradio", "Hugging Face"},
			Features: []string{
				"Warfarin dose prediction using 15+ features",
				"Multiple ML models with 5-fold cross-validation",
				"18% MAE improvement through domain-specific features",
			},
			Category: "AI/ML",
		},
		{
			Name:         "LLM Chatbot",
			Description:  "I built this Go chatbot service to turn my portfolio into an interactive assistant grounded in structured personal knowledge.",
			Backstory:    "Static portfolios felt boring, so I thought, why not let people just ask me things? Building it in Go was a deliberate choice to get more comfortable with the language in a real project.",
			Technologies: []string{"Go", "Gin", "React", "NVIDIA NIM", "OpenRouter", "JWT", "Docker"},
			Features: []string{
				"Conversational portfolio assistant backed by structured personal and project context",
				"Dual-provider setup with NVIDIA NIM primary and OpenRouter fallback",
				"Authentication, guardrails, and rate limiting for production-ready chat APIs",
			},
			Links: map[string]string{
				"Source Code": "https://github.com/chintakjoshi/go-llm-chatbot",
			},
			Category: "AI/ML",
		},
		{
			Name:         "LLM via SMS",
			Description:  "I built this proof of concept so users can chat with an AI assistant over plain SMS through Twilio.",
			Backstory:    "I thought it would be funny to text an AI like it's 2005. Turns out SMS as an interface for LLMs is surprisingly practical for people without reliable internet or smartphones.",
			Technologies: []string{"Python", "Flask", "Twilio", "NVIDIA NIM", "Docker", "ngrok"},
			Features: []string{
				"Twilio webhook flow for inbound and outbound SMS conversations",
				"NVIDIA NIM integration through an OpenAI-compatible client",
				"Containerized local setup with ngrok for webhook exposure",
			},
			Links: map[string]string{
				"Source Code": "https://github.com/chintakjoshi/TxtAI",
			},
			Category: "AI/ML",
		},
		{
			Name:         "authSDK",
			Description:  "I built authSDK as a central authentication platform plus a reusable SDK for downstream services.",
			Backstory:    "After building auth for individual projects multiple times, I got tired of reinventing the wheel. So I built a reusable SDK that any service can just drop in. Laziness is a great motivator.",
			Technologies: []string{"Python", "FastAPI", "PostgreSQL", "Redis", "JWT/JWKS", "OAuth2", "SAML", "Docker"},
			Features: []string{
				"Central auth service for sessions, token issuance, API keys, and lifecycle flows",
				"Reusable SDK for local JWT and API key validation in downstream services",
				"Support for OTP, browser sessions, admin APIs, signing-key rotation, and webhooks",
			},
			Links: map[string]string{
				"Source Code": "https://github.com/chintakjoshi/authSDK",
			},
			Category: "Open Source",
		},
	}
}
