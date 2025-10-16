import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Play, Clock, AlertCircle, CheckCircle } from "lucide-react";
import { runQuery, getShards, QueryResult, Shard } from "@/lib/api";
import { useToast } from "@/hooks/use-toast";

interface QueryRunnerProps {
  projectId: string;
}

export function QueryRunner({ projectId }: QueryRunnerProps) {
  const [query, setQuery] = useState("SELECT * FROM users LIMIT 10;");
  const [selectedShard, setSelectedShard] = useState<string>("all");
  const [result, setResult] = useState<QueryResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [shards, setShards] = useState<Shard[]>([]);
  const { toast } = useToast();

  useEffect(() => {
    const loadShards = async () => {
      try {
        const data = await getShards(projectId);
        setShards(data);
      } catch (error) {
        console.error('Failed to load shards:', error);
      }
    };

    loadShards();
  }, [projectId]);

  const handleRunQuery = async () => {
    if (!query.trim()) {
      toast({
        title: "Error",
        description: "Please enter a SQL query",
        variant: "destructive",
      });
      return;
    }

    setLoading(true);
    try {
      const shardId = selectedShard === "all" ? undefined : selectedShard;
      const queryResult = await runQuery(projectId, query, shardId);
      setResult(queryResult);

      if (queryResult.error) {
        toast({
          title: "Query Error",
          description: queryResult.error,
          variant: "destructive",
        });
      } else {
        toast({
          title: "Success",
          description: `Query executed in ${queryResult.executionTime}s`,
        });
      }
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to execute query",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* Query Input */}
      <Card className="card-gradient">
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="flex items-center space-x-2">
                <Play className="w-5 h-5 text-primary" />
                <span>SQL Query Editor</span>
              </CardTitle>
              <CardDescription>
                Execute SQL queries across your database shards
              </CardDescription>
            </div>
            <div className="flex items-center space-x-3">
              <Select value={selectedShard} onValueChange={setSelectedShard}>
                <SelectTrigger className="w-48">
                  <SelectValue placeholder="Select shard" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Shards</SelectItem>
                  {shards.map((shard) => (
                    <SelectItem key={shard.id} value={shard.id}>
                      {shard.hostname}:{shard.port}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <Button
                onClick={handleRunQuery}
                disabled={loading}
                className="bg-primary hover:bg-primary-glow text-primary-foreground glow-effect"
              >
                <Play className="w-4 h-4 mr-2" />
                {loading ? "Running..." : "Run Query"}
              </Button>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <textarea
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Enter your SQL query here..."
            className="query-editor w-full"
            rows={8}
          />
          <div className="flex items-center justify-between mt-3 text-xs text-muted-foreground">
            <span>Target: {selectedShard === "all" ? "All Shards" : shards.find(s => s.id === selectedShard)?.hostname || "Unknown"}</span>
            <span>{query.length} characters</span>
          </div>
        </CardContent>
      </Card>

      {/* Query Results */}
      {result && (
        <Card className="card-gradient">
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center space-x-2">
                {result.error ? (
                  <AlertCircle className="w-5 h-5 text-destructive" />
                ) : (
                  <CheckCircle className="w-5 h-5 text-success" />
                )}
                <span>Query Results</span>
              </CardTitle>
              <div className="flex items-center space-x-4">
                <Badge variant="outline" className="flex items-center space-x-1">
                  <Clock className="w-3 h-3" />
                  <span>{result.executionTime}s</span>
                </Badge>
                {result.affectedRows !== undefined && (
                  <Badge variant="outline">
                    {result.affectedRows} rows affected
                  </Badge>
                )}
              </div>
            </div>
          </CardHeader>
          <CardContent>
            {result.error ? (
              <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>{result.error}</AlertDescription>
              </Alert>
            ) : result.rows.length > 0 ? (
              <div className="overflow-x-auto">
                <table className="data-table w-full">
                  <thead>
                    <tr>
                      {result.columns.map((column, index) => (
                        <th key={index}>{column}</th>
                      ))}
                    </tr>
                  </thead>
                  <tbody>
                    {result.rows.map((row, rowIndex) => (
                      <tr key={rowIndex}>
                        {row.map((cell, cellIndex) => (
                          <td key={cellIndex} className="font-mono text-sm">
                            {cell?.toString() || <span className="text-muted-foreground">NULL</span>}
                          </td>
                        ))}
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            ) : (
              <div className="text-center py-8 text-muted-foreground">
                <CheckCircle className="w-8 h-8 mx-auto mb-2 text-success" />
                <p>Query executed successfully</p>
                {result.affectedRows !== undefined && (
                  <p className="text-sm">{result.affectedRows} rows affected</p>
                )}
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* Sample Queries */}
      <Card className="card-gradient">
        <CardHeader>
          <CardTitle>Sample Queries</CardTitle>
          <CardDescription>
            Try these example queries to test your sharded database
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {[
              "SELECT COUNT(*) FROM users;",
              "SELECT * FROM orders WHERE status = 'pending';",
              "INSERT INTO users (name, email) VALUES ('John Doe', 'john@example.com');",
              "UPDATE products SET price = price * 1.1 WHERE category = 'electronics';"
            ].map((sampleQuery, index) => (
              <Button
                key={index}
                variant="outline"
                className="justify-start font-mono text-xs h-auto p-3 border-border hover:bg-accent"
                onClick={() => setQuery(sampleQuery)}
              >
                {sampleQuery}
              </Button>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}