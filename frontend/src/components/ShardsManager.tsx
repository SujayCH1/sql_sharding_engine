import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Progress } from "@/components/ui/progress";
import { Plus, Database, Activity, Wifi, WifiOff, Settings, Archive, ArrowLeftRight } from "lucide-react";
import { getShards, addShard, Shard } from "@/lib/api";
import { useToast } from "@/hooks/use-toast";
import { Toaster } from "@/components/ui/toaster";
interface ShardsManagerProps {
  projectId: string;
}

export function ShardsManager({ projectId }: ShardsManagerProps) {
  const [shards, setShards] = useState<Shard[]>([]);
  const [loading, setLoading] = useState(true);
  const [addShardOpen, setAddShardOpen] = useState(false);
  const [newShard, setNewShard] = useState({
    hostname: "",
    port: 5432,
    database: "",
  });
  const { toast } = useToast();

  useEffect(() => {
    loadShards();
  }, [projectId]);

  const loadShards = async () => {
    try {
      const data = await getShards(projectId);
      setShards(data);
    } catch (error) {
      console.error('Failed to load shards:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAddShard = async () => {
    if (!newShard.hostname.trim() || !newShard.database.trim()) {
      toast({
        title: "Error",
        description: "Please fill in all fields",
        variant: "destructive",
      });
      return;
    }

    try {
      const shard = await addShard({
        ...newShard,
        projectId,
      });
      setShards([...shards, shard]);
      setAddShardOpen(false);
      setNewShard({ hostname: "", port: 5432, database: "" });
      toast({
        title: "Success",
        description: "Shard added successfully",
      });
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to add shard",
        variant: "destructive",
      });
    }
  };


  const deleteShard = async (shardId: string) => {
  if (!confirm("Are you sure you want to delete this shard?")) return;

  try {
    // Assume you have a deleteShard API method
    // await deleteShardApi(shardId);
    setShards((prev) => prev.filter((shard) => shard.id !== shardId));
    toast({
      title: "Deleted",
      description: "Shard deleted successfully.",
    });
  } catch (error) {
    toast({
      title: "Error",
      description: "Failed to delete shard.",
      variant: "destructive",
    });
  }
};

  const getStatusIcon = (status: Shard['status']) => {
    switch (status) {
      case 'online':
        return <Wifi className="w-4 h-4 text-success" />;
      case 'offline':
        return <WifiOff className="w-4 h-4 text-destructive" />;
      case 'maintenance':
        return <Settings className="w-4 h-4 text-warning" />;
      default:
        return <WifiOff className="w-4 h-4 text-muted-foreground" />;
    }
  };

  const getStatusBadge = (status: Shard['status']) => {
    const variants = {
      online: 'default',
      offline: 'destructive',
      maintenance: 'secondary',
    } as const;
    
    return (
      <Badge variant={variants[status]} className="flex items-center space-x-1">
        {getStatusIcon(status)}
        <span className="capitalize">{status}</span>
      </Badge>
    );
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };
 

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="animate-pulse">
          <div className="h-8 bg-muted rounded w-1/4 mb-4"></div>
          <div className="space-y-4">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="h-24 bg-muted rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-foreground">Database Shards</h2>
          <p className="text-muted-foreground mt-1">
            Manage your distributed database instances
          </p>
        </div>
        <Dialog open={addShardOpen} onOpenChange={setAddShardOpen}>
          <DialogTrigger asChild>
            <Button className="bg-primary hover:bg-primary-glow text-primary-foreground glow-effect">
              <Plus className="w-4 h-4 mr-2" />
              Add Shard
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-md">
            <DialogHeader>
              <DialogTitle>Add New Shard</DialogTitle>
              <DialogDescription>
                Configure a new database shard for this project
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="hostname">Hostname</Label>
                <Input
                  id="hostname"
                  placeholder="shard1.example.com"
                  value={newShard.hostname}
                  onChange={(e) => setNewShard({ ...newShard, hostname: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="port">Port</Label>
                <Input
                  id="port"
                  type="number"
                  placeholder="5432"
                  value={newShard.port}
                  onChange={(e) => setNewShard({ ...newShard, port: parseInt(e.target.value) || 5432 })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="database">Shard Name</Label>
                <Input
                  id="database"
                  placeholder="project_shard_1"
                  value={newShard.database}
                  onChange={(e) => setNewShard({ ...newShard, database: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="database">Username</Label>
                <Input
                  id="database"
                  placeholder="project_shard_1"
                  value={newShard.database}
                  onChange={(e) => setNewShard({ ...newShard, database: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="database">Password</Label>
                <Input
                  id="database"
                  placeholder="project_shard_1"
                  value={newShard.database}
                  onChange={(e) => setNewShard({ ...newShard, database: e.target.value })}
                />
              </div>
               
              
              <div className="flex space-x-3 pt-4">
                <Button onClick={handleAddShard} className="flex-1">
                  Add Shard
                </Button>
                <Button variant="outline" onClick={() => setAddShardOpen(false)}>
                  Cancel
                </Button>
              </div>
              
            </div>
          </DialogContent>
        </Dialog>
      </div>

      {/* Shards List */}
      {shards.length === 0 ? (
        <Card className="card-gradient">
          <CardHeader className="text-center py-12">
            <Database className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
            <CardTitle>No Shards Configured</CardTitle>
            <CardDescription>
              Add database shards to distribute your data across multiple instances
            </CardDescription>
          </CardHeader>
        </Card>
      ) : (
        <div className="grid gap-6">
          {shards.map((shard) => (
            <Card key={shard.id} className="card-gradient hover:glow-effect transition-smooth">
              <CardHeader>
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-3">
                    <div className="w-10 h-10 bg-primary/20 rounded-lg flex items-center justify-center">
                      <Database className="w-5 h-5 text-primary" />
                    </div>
                    <div>
                      <CardTitle className="text-lg font-semibold">
                        {shard.hostname}:{shard.port}
                      </CardTitle>
                      <CardDescription className="font-mono">
                        {shard.database}
                      </CardDescription>
                    </div>
                  </div>
                  <div className="flex items-center space-x-3">
                    {getStatusBadge(shard.status)}
                    <div className="flex space-x-2">
                      <Button variant="destructive" size="sm" onClick={()=>deleteShard(shard.id)}>
                        <Archive className="w-4 h-4 mr-1" />
                        Delete
                      </Button>
                      {/* <Button variant="outline" size="sm">
                        <ArrowLeftRight className="w-4 h-4 mr-1" />
                        Migrate
                      </Button> */}
                    </div>
                  </div>
                </div>
              </CardHeader>

              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                  {/* Load Progress */}
                  <div className="space-y-2">
                    <div className="flex items-center justify-between text-sm">
                      <span className="text-muted-foreground flex items-center space-x-1">
                        <Activity className="w-3 h-3" />
                        <span>Load</span>
                      </span>
                      <span className="font-medium text-card-foreground">{shard.load}%</span>
                    </div>
                    <Progress value={shard.load} className="h-2" />
                  </div>

                  {/* Connections */}
                  <div className="space-y-2">
                    <div className="flex items-center justify-between text-sm">
                      <span className="text-muted-foreground">Connections</span>
                      <span className="font-medium text-card-foreground">{shard.connections}</span>
                    </div>
                    <Progress value={(shard.connections / 100) * 100} className="h-2" />
                  </div>

                  {/* Last Ping */}
                  <div className="space-y-2">
                    <div className="text-sm">
                      <span className="text-muted-foreground">Last Ping</span>
                      <div className="font-medium text-card-foreground text-xs">
                        {formatDate(shard.lastPing)}
                      </div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}