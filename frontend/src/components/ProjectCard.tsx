import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Database, Calendar, Layers } from "lucide-react";
import { Project } from "@/lib/api";
import { Link } from "react-router-dom";

interface ProjectCardProps {
  project: Project;
}

export function ProjectCard({ project }: ProjectCardProps) {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  return (
    <Card className="card-gradient hover:glow-effect transition-smooth group">
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between">
          <div className="flex items-center space-x-3">
            <div className="w-10 h-10 bg-primary/20 rounded-lg flex items-center justify-center">
              <Database className="w-5 h-5 text-primary" />
            </div>
            <div>
              <CardTitle className="text-lg font-semibold text-card-foreground group-hover:text-primary transition-smooth">
                {project.name}
              </CardTitle>
              <div className="flex items-center space-x-1 mt-1">
                <Calendar className="w-3 h-3 text-muted-foreground" />
                <span className="text-xs text-muted-foreground">
                  Created {formatDate(project.createdAt)}
                </span>
              </div>
            </div>
          </div>
          <Badge variant={project.status === 'active' ? 'default' : 'secondary'}>
            {project.status}
          </Badge>
        </div>
      </CardHeader>

      <CardContent className="pb-4">
        <div className="space-y-2">
          <div className="flex items-center justify-between text-sm">
            <span className="text-muted-foreground flex items-center space-x-1">
              <Layers className="w-3 h-3" />
              <span>Shards</span>
            </span>
            <span className="font-medium text-card-foreground">{project.shardsCount}</span>
          </div>
          <div className="text-xs text-muted-foreground font-mono bg-muted p-2 rounded truncate">
            {project.connectionString}
          </div>
        </div>
      </CardContent>

      <CardFooter>
        <Link to={`/projects/${project.id}`} className="w-full">
          <Button className="w-full bg-primary hover:bg-primary-glow text-primary-foreground transition-smooth">
            Select Project
          </Button>
        </Link>
      </CardFooter>
    </Card>
  );
}