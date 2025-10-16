import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { ArrowLeft, Database, Save } from "lucide-react";
import { addProject } from "@/lib/api";
import { useToast } from "@/hooks/use-toast";

export function NewProjectPage() {
  const [name, setName] = useState("");
  // const [connectionString, setConnectionString] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { toast } = useToast();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()  ) {
      toast({
        title: "Error",
        description: "Please fill in all fields",
        variant: "destructive",
      });
      return;
    }

    setLoading(true);
    try {
      const project = await addProject({
        name: name.trim(),
        // connectionString: connectionString.trim(),
      });
      
      toast({
        title: "Success",
        description: "Project created successfully",
      });
      
      navigate(`/projects/${project.id}`);
    } catch (error) {
      toast({
        title: "Error", 
        description: "Failed to create project",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center space-x-4">
        <Button
          variant="ghost"
          onClick={() => navigate("/")}
          className="text-muted-foreground hover:text-foreground"
        >
          <ArrowLeft className="w-4 h-4 mr-2" />
          Back to Home
        </Button>
        <div>
          <h1 className="text-3xl font-bold text-foreground">Create New Project</h1>
          <p className="text-muted-foreground mt-1">
            Add a new database project to manage sharding
          </p>
        </div>
      </div>

      {/* Form */}
      <div className="max-w-2xl">
        <Card className="card-gradient">
          <CardHeader>
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-primary/20 rounded-lg flex items-center justify-center">
                <Database className="w-5 h-5 text-primary" />
              </div>
              <div>
                <CardTitle>Project Configuration</CardTitle>
                <CardDescription>
                  Configure your database connection and project settings
                </CardDescription>
              </div>
            </div>
          </CardHeader>

          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-6">
              <div className="space-y-2">
                <Label htmlFor="name">Project Name</Label>
                <Input
                  id="name"
                  placeholder="e.g., E-commerce Database"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="bg-input border-border focus:ring-primary"
                />
                <p className="text-xs text-muted-foreground">
                  A descriptive name for your database project
                </p>
              </div>

              {/* <div className="space-y-2">
                <Label htmlFor="connectionString">Database Connection String</Label>
                <Input
                  id="connectionString"
                  placeholder="postgresql://user:password@localhost:5432/database"
                  value={connectionString}
                  onChange={(e) => setConnectionString(e.target.value)}
                  className="bg-input border-border focus:ring-primary font-mono text-sm"
                />
                <p className="text-xs text-muted-foreground">
                  Primary database connection string (PostgreSQL, MySQL, etc.)
                </p>
              </div> */}

              {/* <div className="bg-muted/50 border border-border rounded-lg p-4">
                <h4 className="font-medium text-foreground mb-2">Next Steps</h4>
                <ul className="text-sm text-muted-foreground space-y-1">
                  <li>• Configure sharding strategy in project settings</li>
                  <li>• Add database shards for horizontal scaling</li>
                  <li>• Run queries across distributed data</li>
                </ul>
              </div> */}

              <div className="flex space-x-3 pt-4">
                <Button
                  type="submit"
                  disabled={loading}
                  className="bg-primary hover:bg-primary-glow text-primary-foreground glow-effect"
                >
                  <Save className="w-4 h-4 mr-2" />
                  {loading ? "Creating..." : "Create Project"}
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => navigate("/")}
                  className="border-border hover:bg-accent"
                >
                  Cancel
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}