package service

type ServiceOutput struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	Details     string   `json:"details"`
	Features    []string `json:"features"`
}

var AllServices = []ServiceOutput{
	{ID: 1, Title: "Web Development", Description: "Custom websites and web applications built with the latest technologies.", Icon: "🌐", Details: "Our web development service covers everything from simple static websites to complex web applications. We use modern frameworks like React, Vue, and Angular, combined with powerful backend technologies to deliver fast, scalable, and secure web solutions.", Features: []string{"Responsive design for all devices", "SEO optimization", "Integration with CMS platforms", "E-commerce solutions"}},
	{ID: 2, Title: "Mobile App Development", Description: "Native and cross-platform mobile applications for iOS and Android.", Icon: "📱", Details: "We create engaging and high-performance mobile apps for both iOS and Android platforms. Our team is proficient in native development as well as cross-platform solutions using frameworks like React Native and Flutter.", Features: []string{"Native iOS and Android development", "Cross-platform development", "App Store and Play Store submission", "Maintenance and updates"}},
	{ID: 3, Title: "Desktop App Development", Description: "Powerful desktop applications for Windows, macOS, and Linux.", Icon: "🖥️", Details: "Our desktop application development service focuses on creating robust, efficient, and user-friendly software for Windows, macOS, and Linux. We use technologies like Electron and .NET to build cross-platform applications that provide native-like experiences.", Features: []string{"Cross-platform compatibility", "Integration with system APIs", "Offline functionality", "Automatic updates"}},
	{ID: 4, Title: "Cloud Solutions", Description: "Scalable and secure cloud infrastructure and services.", Icon: "☁️", Details: "We provide comprehensive cloud solutions to help businesses leverage the power of cloud computing. Our services include cloud migration, infrastructure setup, and management across major cloud platforms like AWS, Azure, and Google Cloud.", Features: []string{"Cloud architecture design", "Migration of existing systems to the cloud", "Serverless application development", "24/7 monitoring and support"}},
}

type ServiceUsecase struct{}

func NewServiceUsecase() *ServiceUsecase {
	return &ServiceUsecase{}
}

func (u *ServiceUsecase) GetAllServices() []ServiceOutput {
	return AllServices
}