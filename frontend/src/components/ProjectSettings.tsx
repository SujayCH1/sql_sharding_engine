import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Switch } from "@/components/ui/switch";
import { Separator } from "@/components/ui/separator";
import { Settings, Database, Shield, Zap, Save } from "lucide-react";
import { deleteProject, Project, updateProjectSettings } from "@/lib/api";
import { useToast } from "@/hooks/use-toast";
import { useNavigate } from "react-router-dom";
interface ProjectSettingsProps {
  project: Project;
}

export function ProjectSettings({ project }: ProjectSettingsProps) {
  const navigate = useNavigate();
  const [settings, setSettings] = useState({
    shardingStrategy: "range",
    shardKey: "id",
    replicationFactor: 2,
    autoFailover: true,
    readReplicas: true,
    backupEnabled: true,
    compressionEnabled: false,
    encryptionEnabled: true,
  });
  const [loading, setLoading] = useState(false);
  const { toast } = useToast();

  const handleSave = async () => {
    setLoading(true);
    try {
      await updateProjectSettings(project.id, settings);
      toast({
        title: "Success",
        description: "Project settings updated successfully",
      });
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to update settings",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };


  const handleDeleteProject = async () => {
    const confirmed = window.confirm("Are you sure you want to delete this project?");
    if (!confirmed) return;

    setLoading(true);

    try {
      await deleteProject(project.id);
      toast({
        title: "Deleted",
        description: "Project has been deleted.",
      });
      navigate("/"); 
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to delete project.",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h2 className="text-2xl font-bold text-foreground">Project Settings</h2>
        <p className="text-muted-foreground mt-1">
          Configure database connection and sharding strategies
        </p>
      </div>

      <div className="grid gap-6">
        {/* Database Configuration */}


        {/* <Card className="card-gradient">

          <CardHeader>
            <div className="flex items-center space-x-3">
              <div className="w-8 h-8 bg-primary/20 rounded-lg flex items-center justify-center">
                <Database className="w-4 h-4 text-primary" />
              </div>
              <div>
                <CardTitle>Database Configuration</CardTitle>
                <CardDescription>
                  Primary database connection settings
                </CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="connectionString">Connection String</Label>
              <Input
                id="connectionString"
                value={settings.connectionString}
                onChange={(e) => setSettings({ ...settings, connectionString: e.target.value })}
                className="font-mono text-sm"
              />
              <p className="text-xs text-muted-foreground">
                Primary database connection string
              </p>
            </div>
          </CardContent>
        </Card> */}

        {/* Sharding Strategy */}
        {/* <Card className="card-gradient">
          <CardHeader>
            <div className="flex items-center space-x-3">
              <div className="w-8 h-8 bg-primary/20 rounded-lg flex items-center justify-center">
                <Zap className="w-4 h-4 text-primary" />
              </div>
              <div>
                <CardTitle>Sharding Strategy</CardTitle>
                <CardDescription>
                  Configure how data is distributed across shards
                </CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="shardingStrategy">Strategy</Label>
                <Select
                  value={settings.shardingStrategy}
                  onValueChange={(value) => setSettings({ ...settings, shardingStrategy: value })}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="range">Range-based</SelectItem>
                    <SelectItem value="hash">Hash-based</SelectItem>
                    <SelectItem value="consistent">Consistent Hashing</SelectItem>
                    <SelectItem value="directory">Directory-based</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="shardKey">Shard Key</Label>
                <Input
                  id="shardKey"
                  value={settings.shardKey}
                  onChange={(e) => setSettings({ ...settings, shardKey: e.target.value })}
                  placeholder="e.g., user_id, id"
                />
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="replicationFactor">Replication Factor</Label>
              <Select
                value={settings.replicationFactor.toString()}
                onValueChange={(value) => setSettings({ ...settings, replicationFactor: parseInt(value) })}
              >
                <SelectTrigger className="w-32">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="1">1 (No replication)</SelectItem>
                  <SelectItem value="2">2</SelectItem>
                  <SelectItem value="3">3 (Recommended)</SelectItem>
                  <SelectItem value="4">4</SelectItem>
                  <SelectItem value="5">5</SelectItem>
                </SelectContent>
              </Select>
              <p className="text-xs text-muted-foreground">
                Number of copies of each data partition
              </p>
            </div>
          </CardContent>
        </Card> */}

        {/* High Availability */}
        {/* <Card className="card-gradient">
          <CardHeader>
            <div className="flex items-center space-x-3">
              <div className="w-8 h-8 bg-primary/20 rounded-lg flex items-center justify-center">
                <Shield className="w-4 h-4 text-primary" />
              </div>
              <div>
                <CardTitle>High Availability</CardTitle>
                <CardDescription>
                  Failover and reliability settings
                </CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="autoFailover">Automatic Failover</Label>
                <div className="text-sm text-muted-foreground">
                  Automatically switch to backup shards when primary fails
                </div>
              </div>
              <Switch
                id="autoFailover"
                checked={settings.autoFailover}
                onCheckedChange={(checked) => setSettings({ ...settings, autoFailover: checked })}
              />
            </div>

            <Separator />

            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="readReplicas">Read Replicas</Label>
                <div className="text-sm text-muted-foreground">
                  Enable read-only replicas for improved performance
                </div>
              </div>
              <Switch
                id="readReplicas"
                checked={settings.readReplicas}
                onCheckedChange={(checked) => setSettings({ ...settings, readReplicas: checked })}
              />
            </div>

            <Separator />

            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="backupEnabled">Automated Backups</Label>
                <div className="text-sm text-muted-foreground">
                  Schedule regular backups of all shards
                </div>
              </div>
              <Switch
                id="backupEnabled"
                checked={settings.backupEnabled}
                onCheckedChange={(checked) => setSettings({ ...settings, backupEnabled: checked })}
              />
            </div>
          </CardContent>
        </Card> */}

        {/* Performance & Security */}
        {/* <Card className="card-gradient">
          <CardHeader>
            <div className="flex items-center space-x-3">
              <div className="w-8 h-8 bg-primary/20 rounded-lg flex items-center justify-center">
                <Settings className="w-4 h-4 text-primary" />
              </div>
              <div>
                <CardTitle>Performance & Security</CardTitle>
                <CardDescription>
                  Optimization and security features
                </CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="compressionEnabled">Data Compression</Label>
                <div className="text-sm text-muted-foreground">
                  Compress data to reduce storage requirements
                </div>
              </div>
              <Switch
                id="compressionEnabled"
                checked={settings.compressionEnabled}
                onCheckedChange={(checked) => setSettings({ ...settings, compressionEnabled: checked })}
              />
            </div>

            <Separator />

            <div className="flex items-center justify-between">
              <div className="space-y-0.5">
                <Label htmlFor="encryptionEnabled">Encryption at Rest</Label>
                <div className="text-sm text-muted-foreground">
                  Encrypt stored data for enhanced security
                </div>
              </div>
              <Switch
                id="encryptionEnabled"
                checked={settings.encryptionEnabled}
                onCheckedChange={(checked) => setSettings({ ...settings, encryptionEnabled: checked })}
              />
            </div>
          </CardContent>
        </Card> */}

        {/* Save Button */}
        <div className="flex justify-start">
          <Button
            onClick={handleDeleteProject}
            disabled={loading}
            className="bg-destructive hover:bg-primary-glow text-primary-foreground glow-effect"
          >
            <Save className="w-4 h-4 mr-2" />
            {loading ? "Deleting..." : "Delete Project"}
          </Button>

        </div>



      </div>
    </div>
  );
}