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
			Description:  "Speech Automation Tool designed to automate speech evaluation of language learners and provide instantaneous corrective feedback.",
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
			Name:         "Garbage Collector Ball in Space Game",
			Description:  "3D interactive game where players control a sphere navigating through a galactic environment to collect space debris.",
			Technologies: []string{"Three.js", "WebGL"},
			Features: []string{
				"3D interactive environment",
				"Space debris collection mechanics",
				"Promotes clean virtual universe concept",
			},
			Links: map[string]string{
				"Play Game":   "https://rolling-ball-three-js-chintak-joshis-projects.vercel.app/",
				"Source Code": "https://github.com/chintakjoshi/rolling-ball-three.js",
			},
			Category: "Personal",
		},
		{
			Name:         "Token-based Authentication App",
			Description:  "Authentication application supporting token-based authentication and email verification by OTP.",
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
			Description:  "Web application to help justice-involved job seekers and second-chance employees find employment opportunities.",
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
			Name:         "Online Grocery Store for SLU",
			Description:  "Team lead for project managing grocery stock and facilitating online orders for Saint Louis University store.",
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
			Name:         "E-verification of Documents",
			Description:  "Document validation system that generates unique hex codes and QR codes for identity verification.",
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
			Description:  "Multi-agent system to scan and summarize scientific papers from arXiv and PubMed in real time.",
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
			Name:         "Warfarin Dose Prediction System",
			Description:  "Machine learning system to predict Warfarin dosage using 15+ features, deployed to Gradio via Hugging Face.",
			Technologies: []string{"React", "Flask", "scikit-learn", "Keras", "TensorFlow", "Pandas", "NumPy", "Gradio", "Hugging Face"},
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
			Technologies: []string{"React", "Gin", "Go", "OpenRouter API", "Docker"},
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
