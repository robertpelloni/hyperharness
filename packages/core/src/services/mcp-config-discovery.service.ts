/**
 * @file mcp-config-discovery.service.ts
 * @module packages/core/src/services/mcp-config-discovery.service
 *
 * borg-native entry point for loading MCP server definitions from config.
 * The implementation still lives in the compatibility-named module for now,
 * but active imports can use this generic borg-owned path immediately.
 */

export { getMcpServers } from "./fetch-metamcp.service.js";