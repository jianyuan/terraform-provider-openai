# Import an existing API key for the default project
terraform import openai_project_api_key.example organisation-id/secret-key

# Example
terraform import openai_project_api_key.example org-000000000000000000000000/sk-my-secret-key-xxxxx

# Import an existing API key for a specific project
terraform import openai_project_api_key.example organisation-id/project-id/secret-key

# Example
terraform import openai_project_api_key.example org-000000000000000000000000/proj_000000000000000000000000/sk-my-secret-key-xxxxx
