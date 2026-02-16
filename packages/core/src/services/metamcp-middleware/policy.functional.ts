import { CallToolRequestSchema } from "@modelcontextprotocol/sdk/types.js";
import { CallToolMiddleware } from "./functional-middleware.js";
import { policyService } from "../stubs/policy.service.stub.js";

export interface PolicyMiddlewareOptions {
    enabled: boolean;
}

// @ts-ignore
export const createPolicyMiddleware = (
    options: PolicyMiddlewareOptions,
): CallToolMiddleware => {
    return (next) => async (request, context) => {
        if (!options.enabled) {
            return await next(request, context);
        }

        const { name, _meta } = request.params;

        // Check if a policy ID is attached to the request (e.g. from run_agent)
        // @ts-ignore
        const policyId = _meta?.policyId as string | undefined;

        if (policyId) {
            try {
                const policy = await policyService.getPolicy(policyId);
                if (policy) {
                    const allowed = policyService.evaluateAccess(policy, name);
                    if (!allowed) {
                        throw new Error(`Access denied to tool '${name}' by policy '${policyId}'.`);
                    }
                } else {
                    console.warn(`Policy ID ${policyId} provided but not found.`);
                }
            } catch (error) {
                // Re-throw denial or DB errors
                throw error;
            }
        }

        // Proceed if no policy ID (open access) OR if policy allowed it
        return await next(request, context);
    };
};
