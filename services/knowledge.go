package services

// Project represents a single project with all its details
type Project struct {
	Name         string
	Description  string
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
	ContactInfo  string
}

// GetPersonalInfo returns structured personal information
func GetPersonalInfo() PersonalInfo {
	return PersonalInfo{
		Name:        "Chintak Joshi",
		Title:       "Full-Stack Developer",
		Experience:  "4+ years of experience creating and designing highly scalable, reliable solutions",
		CurrentRole: "Contractor at United Rentals",
		Education: []string{
			"Saint Louis University - Master of Science in Computer Science",
			"Gujarat Technological University - Bachelor's of Engineering in Information Technology",
		},
		Achievements: []string{
			"Proud member of Open Source with SLU, where I led open-source projects across various stacks",
			"Selected for and participated in prestigious national workshops like the IIT Hackathon and Computer Vision",
			"Participated in the Microsoft AI Challenge and successfully completed an introductory course on AI",
		},
		Skills: map[string][]string{
			"Languages":      {"Java", "Python", "JavaScript", "TypeScript", "PHP", "Shell Scripting"},
			"Frameworks":     {"ReactJS", "Spring Boot", "FastAPI", "NodeJS", "Flask", "Django"},
			"AI/ML":          {"TensorFlow", "PyTorch", "Scikit-learn", "Keras", "Pandas", "OpenCV", "Tesseract", "MCP", "LangChain"},
			"Databases":      {"MySQL", "PostgreSQL", "MongoDB", "DynamoDB"},
			"DevOps & Cloud": {"AWS (EC2, S3, RDS, Lambda)", "Elastic", "Docker", "Git", "Jenkins", "Terraform", "CircleCI", "Jira"},
		},
		Links: map[string]string{
			"LinkedIn":  "https://www.linkedin.com/in/chintakjoshi",
			"GitHub":    "https://github.com/chintakjoshi",
			"Portfolio": "https://chintakjoshi.github.io",
		},
		ContactInfo: "To reach out to me, please fill out the contact form on my portfolio. I will reply as soon as possible. Thank you!",
	}
}

// GetProjects returns all actual projects
func GetProjects() []Project {
	return []Project{
		{
			Name:         "iSpraak",
			Description:  "Speech Automation Tool designed to automate speech evaluation of language learners and provide instantaneous corrective feedback.",
			Technologies: []string{"PHP", "ReactJS", "AWS TTS", "Google Auth2.0", "Tailwind", "Docker", "AWS", "S3", "EC2"},
			Features: []string{
				"Automated speech evaluation",
				"Instant corrective feedback",
				"Responsive design with modern UI",
			},
			Links: map[string]string{
				"Live Demo":   "Click Live button in portfolio",
				"Source Code": "Click GitHub in portfolio",
			},
			Category: "Academic/Personal",
		},
		{
			Name:         "Garbage Collector Ball in Space Game",
			Description:  "3D interactive game where players control a sphere navigating through a galactic environment to collect space debris.",
			Technologies: []string{"ThreeJS", "WebGL"},
			Features: []string{
				"3D interactive environment",
				"Space debris collection mechanics",
				"Promotes clean virtual universe concept",
			},
			Links: map[string]string{
				"Play Game":   "Click play button in portfolio",
				"Source Code": "Click GitHub in portfolio",
			},
			Category: "Personal",
		},
		{
			Name:         "Token-based Authentication App",
			Description:  "Authentication application supporting token-based authentication and email verification by OTP.",
			Technologies: []string{"ReactJS", "Spring Boot", "PostgreSQL", "Docker", "Java", "TypeScript"},
			Features: []string{
				"Token-based authentication",
				"Email verification via OTP",
				"Spring Boot backend with React frontend",
			},
			Links: map[string]string{
				"Source Code": "Click GitHub in portfolio",
			},
			Category: "Academic",
		},
		{
			Name:         "Transformative Workforce Academy",
			Description:  "Web application to help justice-involved job seekers and second-chance employees find employment opportunities.",
			Technologies: []string{"ReactJS", "CORS", "NodeJS", "ExpressJS", "PostgreSQL", "Docker"},
			Features: []string{
				"Hiring platform for justice-involved talent",
				"Promotes self-sufficiency and public safety",
				"Open source project for St. Louis community",
			},
			Links: map[string]string{
				"Source Code": "Click GitHub in portfolio",
			},
			Category: "Open Source",
		},
		{
			Name:         "Online Grocery Store for SLU",
			Description:  "Team lead for project managing grocery stock and facilitating online orders for Saint Louis University store.",
			Technologies: []string{"ReactJS", "CORS", "NodeJS", "ExpressJS", "PostgreSQL", "Docker"},
			Features: []string{
				"Grocery stock management",
				"Online ordering system",
				"Team leadership role",
			},
			Links: map[string]string{
				"Source Code": "Click GitHub in portfolio",
			},
			Category: "Academic",
		},
		{
			Name:         "E-verification of Documents",
			Description:  "Document validation system that generates unique hex codes and QR codes for identity verification.",
			Technologies: []string{"Django", "MySQL", "ORM"},
			Features: []string{
				"Document validation against system dataset",
				"Unique hex code generation",
				"QR code for identity verification",
			},
			Links: map[string]string{
				"Source Code": "Click GitHub in portfolio",
			},
			Category: "Academic",
		},
		{
			Name:         "Autonomous Research Assistant for Scientific Literature",
			Description:  "Multi-agent system to scan and summarize scientific papers from arXiv and PubMed in real time.",
			Technologies: []string{"ReactJS", "Flask", "FastAPI", "MCP", "LangChain", "Pinecone"},
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
			Name:         "Warfarin Dose Prediction System",
			Description:  "Machine learning system to predict Warfarin dosage using 15+ features, deployed to Gradio via Hugging Face.",
			Technologies: []string{"ReactJS", "Flask", "Scikit-learn", "Keras", "TensorFlow", "Pandas", "NumPy", "Gradio", "Hugging Face"},
			Features: []string{
				"Warfarin dose prediction using 15+ features",
				"Multiple ML models with 5-fold cross-validation",
				"18% MAE improvement through domain-specific features",
			},
			Category: "AI/ML",
		},
		{
			Name:         "Chatbot using OpenRouter LLM",
			Description:  "Chatbot application leveraging OpenRouter's LLM model's to provide context-aware responses based on structured personal knowledge.",
			Technologies: []string{"ReactJS", "Gin", "Go", "OpenRouter API", "Docker"},
			Features: []string{
				"Context-aware responses using structured knowledge",
				"Integration with OpenRouter's LLM models",
				"Secure authentication and rate limiting",
				"Interactive and user-friendly chat interface on my portfolio",
			},
			Links: map[string]string{
				"Source Code": "https://github.com/chintakjoshi/chintakjoshi-server",
			},
			Category: "AI/ML",
		},
	}
}
