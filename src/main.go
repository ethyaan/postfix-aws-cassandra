package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sigv4-auth-cassandra-gocql-driver-plugin/sigv4"
	"github.com/gocql/gocql"
)

var cassandraSession *gocql.Session
var keyspaceName string
var awsRegion string
var listenPort int = 9999; // derfault port

func main() {

	// check env variables first 
	checkForConfiguration()

	// Initialize Cassandra session
	var err error
	cassandraSession, err = connectToCassandra()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}
	defer cassandraSession.Close()

	// Listen on TCP port 9999
	listenAddress := fmt.Sprintf("127.0.0.1:%d", listenPort)
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("Failed to listen on port 9999: %v", err)
	}
	defer listener.Close()
	log.Printf("Socketmap server listening on %s", listenAddress)

	// Handle incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

/*
**
**
*/
func checkForConfiguration() {

	// Check and assign keyspaceName
	if value, exists := os.LookupEnv("KEYSPACE_NAME"); exists {
		keyspaceName = value
	} else {
		log.Fatalf("Environment variable KEYSPACE_NAME is not set")
	}

	// Check and assign awsRegion
	if value, exists := os.LookupEnv("AWS_REGION"); exists {
		awsRegion = value
	} else {
		log.Fatalf("Environment variable AWS_REGION is not set")
	}

	// Check and assign listenPort, with default fallback
	if value, exists := os.LookupEnv("LISTEN_PORT"); exists {
		port, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("Invalid port number in LISTEN_PORT: %v", err)
		}
		listenPort = port
	}

	// Print values to verify (optional)
	fmt.Printf("KeyspaceName: %s, AwsRegion: %s, ListenPort: %d\n", keyspaceName, awsRegion, listenPort)
}

/* connectToCassandra establishes a connection to AWS Keyspaces using gocql and SigV4 authentication. */
func connectToCassandra() (*gocql.Session, error) {
	
	// Determine contact point using AWS region
	contactPoint := fmt.Sprintf("cassandra.%s.amazonaws.com", awsRegion)
	log.Printf("Using contact point: %s", contactPoint)

	// Configure gocql Cluster
	cluster := gocql.NewCluster(contactPoint)
	cluster.Port = 9142
	cluster.NumConns = 4
	awsAuth := sigv4.NewAwsAuthenticator()
	cluster.Authenticator = awsAuth

	// TLS configuration for Keyspaces
	cluster.SslOpts = &gocql.SslOptions{
		CaPath:                 "certs/sf-class2-root.crt",
		EnableHostVerification: false,
	}

	cluster.Consistency = gocql.LocalQuorum
	cluster.DisableInitialHostLookup = false

	return cluster.CreateSession()
}

//  handles incoming connections and processes queries.
func handleConnection(conn net.Conn) {
    defer conn.Close()
    reader := bufio.NewReader(conn)
    writer := bufio.NewWriter(conn)

    for {
        // Read incoming request
        line, err := reader.ReadString('\n')
        if err != nil {
            log.Printf("Error reading from connection: %v", err)
            return
        }
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }

        // Process request
        parts := strings.SplitN(line, " ", 2)
        if len(parts) != 2 {
            log.Printf("Invalid request: %s", line)
            continue
        }
        tableName := parts[0]
        key := parts[1]

        // Respond to the query based on the table
        var response string
        switch tableName {
        case "domains":
            response = queryDomains(key)
        case "users":
            response = queryUsers(key)
        case "relay_domains":
            response = queryRelayDomains(key)
        case "virtual_aliases":
            response = queryVirtualAliases(key)
        case "transport_maps":
            response = queryTransportMaps(key)
        case "access_maps":
            response = queryAccessMaps(key)
        default:
            response = "NOTFOUND"
        }

        // Write the response back to the client
        writer.WriteString(response + "\n")
        writer.Flush()
    }
}

func querySingleResult(tableName, selectColumn, whereColumn, key string, result interface{}) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    query := fmt.Sprintf("SELECT %s FROM %s.%s WHERE %s = ?", selectColumn, keyspaceName, tableName, whereColumn)

    // Execute the query
    err := cassandraSession.Query(query, key).WithContext(ctx).Scan(result)
    if err != nil {
        if err == gocql.ErrNotFound {
            return gocql.ErrNotFound
        }
        log.Printf("Error executing query on table %s: %v", tableName, err)
        return err
    }

    return nil
}

func queryDomains(domain string) string {
    var active bool
    err := querySingleResult("domains", "active", "domain", domain, &active)
    if err != nil {
        if err == gocql.ErrNotFound {
            return "NOTFOUND"
        }
        log.Printf("Error querying domains: %v", err)
        return "TEMPFAIL"
    }
    if active {
        return "OK"
    }
    return "NOTFOUND"
}

func queryUsers(email string) string {
    var active bool
    err := querySingleResult("users", "active", "email", email, &active)
    if err != nil {
        if err == gocql.ErrNotFound {
            return "NOTFOUND"
        }
        log.Printf("Error querying users: %v", err)
        return "TEMPFAIL"
    }
    if active {
        return "OK"
    }
    return "NOTFOUND"
}

func queryRelayDomains(domain string) string {
    var active bool
    err := querySingleResult("relay_domains", "active", "domain", domain, &active)
    if err != nil {
        if err == gocql.ErrNotFound {
            return "NOTFOUND"
        }
        log.Printf("Error querying relay_domains: %v", err)
        return "TEMPFAIL"
    }
    if active {
        return "OK"
    }
    return "NOTFOUND"
}

func queryVirtualAliases(alias string) string {
    var destination string
    err := querySingleResult("virtual_aliases", "destination", "alias", alias, &destination)
    if err != nil {
        if err == gocql.ErrNotFound {
            return "NOTFOUND"
        }
        log.Printf("Error querying virtual_aliases: %v", err)
        return "TEMPFAIL"
    }
    return fmt.Sprintf("OK %s", destination)
}

func queryTransportMaps(address string) string {
    var transport string
    err := querySingleResult("transport_maps", "transport", "address", address, &transport)
    if err != nil {
        if err == gocql.ErrNotFound {
            return "NOTFOUND"
        }
        log.Printf("Error querying transport_maps: %v", err)
        return "TEMPFAIL"
    }
    return fmt.Sprintf("OK %s", transport)
}

func queryAccessMaps(sender string) string {
    var action string
    err := querySingleResult("access_maps", "action", "sender", sender, &action)
    if err != nil {
        if err == gocql.ErrNotFound {
            return "DUNNO" // Default action if not found
        }
        log.Printf("Error querying access_maps: %v", err)
        return "DEFER_IF_PERMIT Service unavailable"
    }
    return action
}