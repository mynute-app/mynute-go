// Atlas configuration for mynute-go
// Defines data sources and migration environments

variable "db_url" {
  type    = string
  default = getenv("DATABASE_URL")
}

variable "db_host" {
  type    = string
  default = getenv("DB_HOST")
}

variable "db_port" {
  type    = string
  default = getenv("DB_PORT")
}

variable "db_user" {
  type    = string
  default = getenv("DB_USER")
}

variable "db_password" {
  type    = string
  default = getenv("DB_PASSWORD")
}

variable "db_name" {
  type    = string
  default = getenv("DB_NAME")
}

// Construct database URL from components if DATABASE_URL is not set
locals {
  db_url = var.db_url != "" ? var.db_url : "postgres://${var.db_user}:${var.db_password}@${var.db_host}:${var.db_port}/${var.db_name}?sslmode=disable"
}

// Data source for current database state
data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./cmd/atlas-loader",
  ]
}

// Default environment
env "dev" {
  src = data.external_schema.gorm.url
  url = local.db_url
  
  dev = "docker://postgres/15/dev?search_path=public"
  
  migration {
    dir = "file://migrations"
  }
  
  diff {
    skip {
      drop_schema = true
      drop_table  = true
    }
  }
}

env "prod" {
  src = data.external_schema.gorm.url
  url = local.db_url
  
  migration {
    dir = "file://migrations"
  }
  
  diff {
    skip {
      drop_schema = true
      drop_table  = true
    }
  }
}
