<?php
/**
 * Sailtix API Documentation - Server-Side Authentication Handler
 * This file handles server-side authentication and access control
 */

// Prevent direct access to this file
if (!isset($_SERVER['HTTP_REFERER']) || strpos($_SERVER['HTTP_REFERER'], $_SERVER['HTTP_HOST']) === false) {
    http_response_code(403);
    exit('Access Denied');
}

// Security headers
header('X-Content-Type-Options: nosniff');
header('X-Frame-Options: DENY');
header('X-XSS-Protection: 1; mode=block');
header('Referrer-Policy: strict-origin-when-cross-origin');

// Authorized agents configuration (in production, this should be in a secure database)
$AUTHORIZED_AGENTS = [
    'sailtix-agent-001' => [
        'access_key' => 'access-key-2024-sailtix-secure',
        'name' => 'Primary Sailtix Agent',
        'role' => 'production',
        'permissions' => ['read', 'test']
    ],
    'sailtix-agent-002' => [
        'access_key' => 'access-key-2024-sailtix-secure',
        'name' => 'Secondary Sailtix Agent',
        'role' => 'production',
        'permissions' => ['read', 'test']
    ],
    'sailtix-dev-001' => [
        'access_key' => 'dev-access-key-2024',
        'name' => 'Development Agent',
        'role' => 'development',
        'permissions' => ['read', 'test', 'debug']
    ],
    'sailtix-admin' => [
        'access_key' => 'admin-access-key-2024-secure',
        'name' => 'Administrator',
        'role' => 'admin',
        'permissions' => ['read', 'test', 'debug', 'admin']
    ]
];

// Security settings
$SESSION_DURATION = 24 * 60 * 60; // 24 hours
$MAX_LOGIN_ATTEMPTS = 5;
$LOCKOUT_DURATION = 15 * 60; // 15 minutes

// Rate limiting
session_start();
$current_time = time();

// Check for brute force attempts
if (isset($_SESSION['login_attempts']) && $_SESSION['login_attempts']['count'] >= $MAX_LOGIN_ATTEMPTS) {
    $time_diff = $current_time - $_SESSION['login_attempts']['time'];
    if ($time_diff < $LOCKOUT_DURATION) {
        http_response_code(429);
        exit(json_encode([
            'error' => 'Too many login attempts. Please try again in ' . ceil(($LOCKOUT_DURATION - $time_diff) / 60) . ' minutes.'
        ]));
    } else {
        // Reset attempts after lockout period
        unset($_SESSION['login_attempts']);
    }
}

// Handle authentication request
if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $input = json_decode(file_get_contents('php://input'), true);
    
    if (!$input) {
        http_response_code(400);
        exit(json_encode(['error' => 'Invalid request format']));
    }
    
    $agent_id = trim($input['agentId'] ?? '');
    $access_key = trim($input['accessKey'] ?? '');
    
    // Validate input
    if (empty($agent_id) || empty($access_key)) {
        http_response_code(400);
        exit(json_encode(['error' => 'Agent ID and Access Key are required']));
    }
    
    // Check credentials
    if (isset($AUTHORIZED_AGENTS[$agent_id]) && $AUTHORIZED_AGENTS[$agent_id]['access_key'] === $access_key) {
        // Authentication successful
        unset($_SESSION['login_attempts']); // Reset attempts
        
        // Generate secure session token
        $session_token = bin2hex(random_bytes(32));
        $expires = $current_time + $SESSION_DURATION;
        
        // Store session data (in production, use Redis or database)
        $_SESSION['sailtix_auth'] = [
            'agent_id' => $agent_id,
            'agent_name' => $AUTHORIZED_AGENTS[$agent_id]['name'],
            'role' => $AUTHORIZED_AGENTS[$agent_id]['role'],
            'permissions' => $AUTHORIZED_AGENTS[$agent_id]['permissions'],
            'token' => $session_token,
            'expires' => $expires,
            'ip' => $_SERVER['REMOTE_ADDR'],
            'user_agent' => $_SERVER['HTTP_USER_AGENT']
        ];
        
        // Set secure cookie
        setcookie('sailtix_auth', $session_token, [
            'expires' => $expires,
            'path' => '/',
            'domain' => '',
            'secure' => true, // Requires HTTPS
            'httponly' => true,
            'samesite' => 'Strict'
        ]);
        
        // Log successful login (in production, log to secure log file)
        error_log("Successful login: Agent ID: $agent_id, IP: " . $_SERVER['REMOTE_ADDR'] . ", Time: " . date('Y-m-d H:i:s'));
        
        http_response_code(200);
        exit(json_encode([
            'success' => true,
            'message' => 'Authentication successful',
            'agent' => [
                'id' => $agent_id,
                'name' => $AUTHORIZED_AGENTS[$agent_id]['name'],
                'role' => $AUTHORIZED_AGENTS[$agent_id]['role']
            ]
        ]));
        
    } else {
        // Authentication failed
        if (!isset($_SESSION['login_attempts'])) {
            $_SESSION['login_attempts'] = ['count' => 0, 'time' => $current_time];
        }
        $_SESSION['login_attempts']['count']++;
        $_SESSION['login_attempts']['time'] = $current_time;
        
        // Log failed attempt (in production, log to secure log file)
        error_log("Failed login attempt: Agent ID: $agent_id, IP: " . $_SERVER['REMOTE_ADDR'] . ", Time: " . date('Y-m-d H:i:s'));
        
        http_response_code(401);
        exit(json_encode(['error' => 'Invalid Agent ID or Access Key']));
    }
}

// Handle logout
if ($_SERVER['REQUEST_METHOD'] === 'DELETE') {
    // Clear session
    unset($_SESSION['sailtix_auth']);
    session_destroy();
    
    // Clear cookie
    setcookie('sailtix_auth', '', [
        'expires' => time() - 3600,
        'path' => '/',
        'domain' => '',
        'secure' => true,
        'httponly' => true,
        'samesite' => 'Strict'
    ]);
    
    http_response_code(200);
    exit(json_encode(['success' => true, 'message' => 'Logged out successfully']));
}

// Handle API specification requests
if ($_SERVER['REQUEST_METHOD'] === 'GET') {
    // Check if user is authenticated
    if (!isset($_SESSION['sailtix_auth']) || $_SESSION['sailtix_auth']['expires'] < $current_time) {
        http_response_code(401);
        exit(json_encode(['error' => 'Authentication required']));
    }
    
    // Validate session token
    if (!isset($_COOKIE['sailtix_auth']) || $_COOKIE['sailtix_auth'] !== $_SESSION['sailtix_auth']['token']) {
        http_response_code(401);
        exit(json_encode(['error' => 'Invalid session']));
    }
    
    // Serve API specification
    $requested_file = $_GET['file'] ?? '';
    $allowed_files = ['openapi.yaml', 'openapi-v2.yaml'];
    
    if (!in_array($requested_file, $allowed_files)) {
        http_response_code(404);
        exit(json_encode(['error' => 'File not found']));
    }
    
    $file_path = __DIR__ . '/' . $requested_file;
    if (!file_exists($file_path)) {
        http_response_code(404);
        exit(json_encode(['error' => 'File not found']));
    }
    
    // Set appropriate headers for YAML file
    header('Content-Type: application/x-yaml');
    header('Content-Disposition: inline; filename="' . $requested_file . '"');
    header('Cache-Control: no-cache, no-store, must-revalidate');
    header('Pragma: no-cache');
    header('Expires: 0');
    
    // Output file content
    readfile($file_path);
    exit;
}

// Invalid request method
http_response_code(405);
exit(json_encode(['error' => 'Method not allowed']));
?> 