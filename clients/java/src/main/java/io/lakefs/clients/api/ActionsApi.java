/*
 * lakeFS API
 * lakeFS HTTP API
 *
 * The version of the OpenAPI document: 0.1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package io.lakefs.clients.api;

import io.lakefs.clients.api.ApiCallback;
import io.lakefs.clients.api.ApiClient;
import io.lakefs.clients.api.ApiException;
import io.lakefs.clients.api.ApiResponse;
import io.lakefs.clients.api.Configuration;
import io.lakefs.clients.api.Pair;
import io.lakefs.clients.api.ProgressRequestBody;
import io.lakefs.clients.api.ProgressResponseBody;

import com.google.gson.reflect.TypeToken;

import java.io.IOException;


import io.lakefs.clients.api.model.ActionRun;
import io.lakefs.clients.api.model.ActionRunList;
import io.lakefs.clients.api.model.Error;
import java.io.File;
import io.lakefs.clients.api.model.HookRunList;

import java.lang.reflect.Type;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class ActionsApi {
    private ApiClient localVarApiClient;

    public ActionsApi() {
        this(Configuration.getDefaultApiClient());
    }

    public ActionsApi(ApiClient apiClient) {
        this.localVarApiClient = apiClient;
    }

    public ApiClient getApiClient() {
        return localVarApiClient;
    }

    public void setApiClient(ApiClient apiClient) {
        this.localVarApiClient = apiClient;
    }

    /**
     * Build call for getRun
     * @param repository  (required)
     * @param runId  (required)
     * @param _callback Callback for upload/download progress
     * @return Call to execute
     * @throws ApiException If fail to serialize the request body object
     * @http.response.details
     <table summary="Response Details" border="1">
        <tr><td> Status Code </td><td> Description </td><td> Response Headers </td></tr>
        <tr><td> 200 </td><td> action run result </td><td>  -  </td></tr>
        <tr><td> 401 </td><td> Unauthorized </td><td>  -  </td></tr>
        <tr><td> 404 </td><td> Resource Not Found </td><td>  -  </td></tr>
        <tr><td> 0 </td><td> Internal Server Error </td><td>  -  </td></tr>
     </table>
     */
    public okhttp3.Call getRunCall(String repository, String runId, final ApiCallback _callback) throws ApiException {
        Object localVarPostBody = null;

        // create path and map variables
        String localVarPath = "/repositories/{repository}/actions/runs/{run_id}"
            .replaceAll("\\{" + "repository" + "\\}", localVarApiClient.escapeString(repository.toString()))
            .replaceAll("\\{" + "run_id" + "\\}", localVarApiClient.escapeString(runId.toString()));

        List<Pair> localVarQueryParams = new ArrayList<Pair>();
        List<Pair> localVarCollectionQueryParams = new ArrayList<Pair>();
        Map<String, String> localVarHeaderParams = new HashMap<String, String>();
        Map<String, String> localVarCookieParams = new HashMap<String, String>();
        Map<String, Object> localVarFormParams = new HashMap<String, Object>();

        final String[] localVarAccepts = {
            "application/json"
        };
        final String localVarAccept = localVarApiClient.selectHeaderAccept(localVarAccepts);
        if (localVarAccept != null) {
            localVarHeaderParams.put("Accept", localVarAccept);
        }

        final String[] localVarContentTypes = {
            
        };
        final String localVarContentType = localVarApiClient.selectHeaderContentType(localVarContentTypes);
        localVarHeaderParams.put("Content-Type", localVarContentType);

        String[] localVarAuthNames = new String[] { "basic_auth", "cookie_auth", "jwt_token", "oidc_auth" };
        return localVarApiClient.buildCall(localVarPath, "GET", localVarQueryParams, localVarCollectionQueryParams, localVarPostBody, localVarHeaderParams, localVarCookieParams, localVarFormParams, localVarAuthNames, _callback);
    }

    @SuppressWarnings("rawtypes")
    private okhttp3.Call getRunValidateBeforeCall(String repository, String runId, final ApiCallback _callback) throws ApiException {
        
        // verify the required parameter 'repository' is set
        if (repository == null) {
            throw new ApiException("Missing the required parameter 'repository' when calling getRun(Async)");
        }
        
        // verify the required parameter 'runId' is set
        if (runId == null) {
            throw new ApiException("Missing the required parameter 'runId' when calling getRun(Async)");
        }
        

        okhttp3.Call localVarCall = getRunCall(repository, runId, _callback);
        return localVarCall;

    }

    /**
     * get a run
     * 
     * @param repository  (required)
     * @param runId  (required)
     * @return ActionRun
     * @throws ApiException If fail to call the API, e.g. server error or cannot deserialize the response body
     * @http.response.details
     <table summary="Response Details" border="1">
        <tr><td> Status Code </td><td> Description </td><td> Response Headers </td></tr>
        <tr><td> 200 </td><td> action run result </td><td>  -  </td></tr>
        <tr><td> 401 </td><td> Unauthorized </td><td>  -  </td></tr>
        <tr><td> 404 </td><td> Resource Not Found </td><td>  -  </td></tr>
        <tr><td> 0 </td><td> Internal Server Error </td><td>  -  </td></tr>
     </table>
     */
    public ActionRun getRun(String repository, String runId) throws ApiException {
        ApiResponse<ActionRun> localVarResp = getRunWithHttpInfo(repository, runId);
        return localVarResp.getData();
    }

    /**
     * get a run
     * 
     * @param repository  (required)
     * @param runId  (required)
     * @return ApiResponse&lt;ActionRun&gt;
     * @throws ApiException If fail to call the API, e.g. server error or cannot deserialize the response body
     * @http.response.details
     <table summary="Response Details" border="1">
        <tr><td> Status Code </td><td> Description </td><td> Response Headers </td></tr>
        <tr><td> 200 </td><td> action run result </td><td>  -  </td></tr>
        <tr><td> 401 </td><td> Unauthorized </td><td>  -  </td></tr>
        <tr><td> 404 </td><td> Resource Not Found </td><td>  -  </td></tr>
        <tr><td> 0 </td><td> Internal Server Error </td><td>  -  </td></tr>
     </table>
     */
    public ApiResponse<ActionRun> getRunWithHttpInfo(String repository, String runId) throws ApiException {
        okhttp3.Call localVarCall = getRunValidateBeforeCall(repository, runId, null);
        Type localVarReturnType = new TypeToken<ActionRun>(){}.getType();
        return localVarApiClient.execute(localVarCall, localVarReturnType);
    }

    /**
     * get a run (asynchronously)
     * 
     * @param repository  (required)
     * @param runId  (required)
     * @param _callback The callback to be executed when the API call finishes
     * @return The request call
     * @throws ApiException If fail to process the API call, e.g. serializing the request body object
     * @http.response.details
     <table summary="Response Details" border="1">
        <tr><td> Status Code </td><td> Description </td><td> Response Headers </td></tr>
        <tr><td> 200 </td><td> action run result </td><td>  -  </td></tr>
        <tr><td> 401 </td><td> Unauthorized </td><td>  -  </td></tr>
        <tr><td> 404 </td><td> Resource Not Found </td><td>  -  </td></tr>
        <tr><td> 0 </td><td> Internal Server Error </td><td>  -  </td></tr>
     </table>
     */
    public okhttp3.Call getRunAsync(String repository, String runId, final ApiCallback<ActionRun> _callback) throws ApiException {

        okhttp3.Call localVarCall = getRunValidateBeforeCall(repository, runId, _callback);
        Type localVarReturnType = new TypeToken<ActionRun>(){}.getType();
        localVarApiClient.executeAsync(localVarCall, localVarReturnType, _callback);
        return localVarCall;
    }
    /**
     * Build call for getRunHookOutput
     * @param repository  (required)
     * @param runId  (required)
     * @param hookRunId  (required)
     * @param _callback Callback for upload/download progress
     * @return Call to execute
     * @throws ApiException If fail to serialize the request body object
     * @http.response.details
     <table summary="Response Details" border="1">
        <tr><td> Status Code </td><td> Description </td><td> Response Headers </td></tr>
        <tr><td> 200 </td><td> run hook output </td><td>  -  </td></tr>
        <tr><td> 401 </td><td> Unauthorized </td><td>  -  </td></tr>
        <tr><td> 404 </td><td> Resource Not Found </td><td>  -  </td></tr>
        <tr><td> 0 </td><td> Internal Server Error </td><td>  -  </td></tr>
     </table>
     */
    public okhttp3.Call getRunHookOutputCall(String repository, String runId, String hookRunId, final ApiCallback _callback) throws ApiException {
        Object localVarPostBody = null;

        // create path and map variables
        String localVarPath = "/repositories/{repository}/actions/runs/{run_id}/hooks/{hook_run_id}/output"
            .replaceAll("\\{" + "repository" + "\\}", localVarApiClient.escapeString(repository.toString()))
            .replaceAll("\\{" + "run_id" + "\\}", localVarApiClient.escapeString(runId.toString()))
            .replaceAll("\\{" + "hook_run_id" + "\\}", localVarApiClient.escapeString(hookRunId.toString()));

        List<Pair> localVarQueryParams = new ArrayList<Pair>();
        List<Pair> localVarCollectionQueryParams = new ArrayList<Pair>();
        Map<String, String> localVarHeaderParams = new HashMap<String, String>();
        Map<String, String> localVarCookieParams = new HashMap<String, String>();
        Map<String, Object> localVarFormParams = new HashMap<String, Object>();

        final String[] localVarAccepts = {
            "application/octet-stream", "application/json"
        };
        final String localVarAccept = localVarApiClient.selectHeaderAccept(localVarAccepts);
        if (localVarAccept != null) {
            localVarHeaderParams.put("Accept", localVarAccept);
        }

        final String[] localVarContentTypes = {
            
        };
        final String localVarContentType = localVarApiClient.selectHeaderContentType(localVarContentTypes);
        localVarHeaderParams.put("Content-Type", localVarContentType);

        String[] localVarAuthNames = new String[] { "basic_auth", "cookie_auth", "jwt_token", "oidc_auth" };
        return localVarApiClient.buildCall(localVarPath, "GET", localVarQueryParams, localVarCollectionQueryParams, localVarPostBody, localVarHeaderParams, localVarCookieParams, localVarFormParams, localVarAuthNames, _callback);
    }

    @SuppressWarnings("rawtypes")
    private okhttp3.Call getRunHookOutputValidateBeforeCall(String repository, String runId, String hookRunId, final ApiCallback _callback) throws ApiException {
        
        // verify the required parameter 'repository' is set
        if (repository == null) {
            throw new ApiException("Missing the required parameter 'repository' when calling getRunHookOutput(Async)");
        }
        
        // verify the required parameter 'runId' is set
        if (runId == null) {
            throw new ApiException("Missing the required parameter 'runId' when calling getRunHookOutput(Async)");
        }
        
        // verify the required parameter 'hookRunId' is set
        if (hookRunId == null) {
            throw new ApiException("Missing the required parameter 'hookRunId' when calling getRunHookOutput(Async)");
        }
        

        okhttp3.Call localVarCall = getRunHookOutputCall(repository, runId, hookRunId, _callback);
        return localVarCall;

    }

    /**
     * get run hook output
     * 
     * @param repository  (required)
     * @param runId  (required)
     * @param hookRunId  (required)
     * @return File
     * @throws ApiException If fail to call the API, e.g. server error or cannot deserialize the response body
     * @http.response.details
     <table summary="Response Details" border="1">
        <tr><td> Status Code </td><td> Description </td><td> Response Headers </td></tr>
        <tr><td> 200 </td><td> run hook output </td><td>  -  </td></tr>
        <tr><td> 401 </td><td> Unauthorized </td><td>  -  </td></tr>
        <tr><td> 404 </td><td> Resource Not Found </td><td>  -  </td></tr>
        <tr><td> 0 </td><td> Internal Server Error </td><td>  -  </td></tr>
     </table>
     */
    public File getRunHookOutput(String repository, String runId, String hookRunId) throws ApiException {
        ApiResponse<File> localVarResp = getRunHookOutputWithHttpInfo(repository, runId, hookRunId);
        return localVarResp.getData();
    }

    /**
     * get run hook output
     * 
     * @param repository  (required)
     * @param runId  (required)
     * @param hookRunId  (required)
     * @return ApiResponse&lt;File&gt;
     * @throws ApiException If fail to call the API, e.g. server error or cannot deserialize the response body
     * @http.response.details
     <table summary="Response Details" border="1">
        <tr><td> Status Code </td><td> Description </td><td> Response Headers </td></tr>
        <tr><td> 200 </td><td> run hook output </td><td>  -  </td></tr>
        <tr><td> 401 </td><td> Unauthorized </td><td>  -  </td></tr>
        <tr><td> 404 </td><td> Resource Not Found </td><td>  -  </td></tr>
        <tr><td> 0 </td><td> Internal Server Error </td><td>  -  </td></tr>
     </table>
     */
    public ApiResponse<File> getRunHookOutputWithHttpInfo(String repository, String runId, String hookRunId) throws ApiException {
        okhttp3.Call localVarCall = getRunHookOutputValidateBeforeCall(repository, runId, hookRunId, null);
        Type localVarReturnType = new TypeToken<File>(){}.getType();
        return localVarApiClient.execute(localVarCall, localVarReturnType);
    }

    /**
     * get run hook output (asynchronously)
     * 
     * @param repository  (required)
     * @param runId  (required)
     * @param hookRunId  (required)
     * @param _callback The callback to be executed when the API call finishes
     * @return The request call
     * @throws ApiException If fail to process the API call, e.g. serializing the request body object
     * @http.response.details
     <table summary="Response Details" border="1">
        <tr><td> Status Code </td><td> Description </td><td> Response Headers </td></tr>
        <tr><td> 200 </td><td> run hook output </td><td>  -  </td></tr>
        <tr><td> 401 </td><td> Unauthorized </td><td>  -  </td></tr>
        <tr><td> 404 </td><td> Resource Not Found </td><td>  -  </td></tr>
        <tr><td> 0 </td><td> Internal Server Error </td><td>  -  </td></tr>
     </table>
     */
    public okhttp3.Call getRunHookOutputAsync(String repository, String runId, String hookRunId, final ApiCallback<File> _callback) throws ApiException {

        okhttp3.Call localVarCall = getRunHookOutputValidateBeforeCall(repository, runId, hookRunId, _callback);
        Type localVarReturnType = new TypeToken<File>(){}.getType();
        localVarApiClient.executeAsync(localVarCall, localVarReturnType, _callback);
        return localVarCall;
    }
    /**
     * Build call for listRepositoryRuns
     * @param repository  (required)
     * @param after return items after this value (optional)
     * @param amount how many items to return (optional, default to 100)
     * @param branch  (optional)
     * @param commit  (optional)
     * @param _callback Callback for upload/download progress
     * @return Call to execute
     * @throws ApiException If fail to serialize the request body object
     * @http.response.details
     <table summary="Response Details" border="1">
        <tr><td> Status Code </td><td> Description </td><td> Response Headers </td></tr>
        <tr><td> 200 </td><td> list action runs </td><td>  -  </td></tr>
        <tr><td> 401 </td><td> Unauthorized </td><td>  -  </td></tr>
        <tr><td> 404 </td><td> Resource Not Found </td><td>  -  </td></tr>
        <tr><td> 0 </td><td> Internal Server Error </td><td>  -  </td></tr>
     </table>
     */
    public okhttp3.Call listRepositoryRunsCall(String repository, String after, Integer amount, String branch, String commit, final ApiCallback _callback) throws ApiException {
        Object localVarPostBody = null;

        // create path and map variables
        String localVarPath = "/repositories/{repository}/actions/runs"
            .replaceAll("\\{" + "repository" + "\\}", localVarApiClient.escapeString(repository.toString()));

        List<Pair> localVarQueryParams = new ArrayList<Pair>();
        List<Pair> localVarCollectionQueryParams = new ArrayList<Pair>();
        Map<String, String> localVarHeaderParams = new HashMap<String, String>();
        Map<String, String> localVarCookieParams = new HashMap<String, String>();
        Map<String, Object> localVarFormParams = new HashMap<String, Object>();

        if (after != null) {
            localVarQueryParams.addAll(localVarApiClient.parameterToPair("after", after));
        }

        if (amount != null) {
            localVarQueryParams.addAll(localVarApiClient.parameterToPair("amount", amount));
        }

        if (branch != null) {
            localVarQueryParams.addAll(localVarApiClient.parameterToPair("branch", branch));
        }

        if (commit != null) {
            localVarQueryParams.addAll(localVarApiClient.parameterToPair("commit", commit));
        }

        final String[] localVarAccepts = {
            "application/json"
        };
        final String localVarAccept = localVarApiClient.selectHeaderAccept(localVarAccepts);
        if (localVarAccept != null) {
            localVarHeaderParams.put("Accept", localVarAccept);
        }

        final String[] localVarContentTypes = {
            
        };
        final String localVarContentType = localVarApiClient.selectHeaderContentType(localVarContentTypes);
        localVarHeaderParams.put("Content-Type", localVarContentType);

        String[] localVarAuthNames = new String[] { "basic_auth", "cookie_auth", "jwt_token", "oidc_auth" };
        return localVarApiClient.buildCall(localVarPath, "GET", localVarQueryParams, localVarCollectionQueryParams, localVarPostBody, localVarHeaderParams, localVarCookieParams, localVarFormParams, localVarAuthNames, _callback);
    }

    @SuppressWarnings("rawtypes")
    private okhttp3.Call listRepositoryRunsValidateBeforeCall(String repository, String after, Integer amount, String branch, String commit, final ApiCallback _callback) throws ApiException {
        
        // verify the required parameter 'repository' is set
        if (repository == null) {
            throw new ApiException("Missing the required parameter 'repository' when calling listRepositoryRuns(Async)");
        }
        

        okhttp3.Call localVarCall = listRepositoryRunsCall(repository, after, amount, branch, commit, _callback);
        return localVarCall;

    }

    /**
     * list runs
     * 
     * @param repository  (required)
     * @param after return items after this value (optional)
     * @param amount how many items to return (optional, default to 100)
     * @param branch  (optional)
     * @param commit  (optional)
     * @return ActionRunList
     * @throws ApiException If fail to call the API, e.g. server error or cannot deserialize the response body
     * @http.response.details
     <table summary="Response Details" border="1">
        <tr><td> Status Code </td><td> Description </td><td> Response Headers </td></tr>
        <tr><td> 200 </td><td> list action runs </td><td>  -  </td></tr>
        <tr><td> 401 </td><td> Unauthorized </td><td>  -  </td></tr>
        <tr><td> 404 </td><td> Resource Not Found </td><td>  -  </td></tr>
        <tr><td> 0 </td><td> Internal Server Error </td><td>  -  </td></tr>
     </table>
     */
    public ActionRunList listRepositoryRuns(String repository, String after, Integer amount, String branch, String commit) throws ApiException {
        ApiResponse<ActionRunList> localVarResp = listRepositoryRunsWithHttpInfo(repository, after, amount, branch, commit);
        return localVarResp.getData();
    }

    /**
     * list runs
     * 
     * @param repository  (required)
     * @param after return items after this value (optional)
     * @param amount how many items to return (optional, default to 100)
     * @param branch  (optional)
     * @param commit  (optional)
     * @return ApiResponse&lt;ActionRunList&gt;
     * @throws ApiException If fail to call the API, e.g. server error or cannot deserialize the response body
     * @http.response.details
     <table summary="Response Details" border="1">
        <tr><td> Status Code </td><td> Description </td><td> Response Headers </td></tr>
        <tr><td> 200 </td><td> list action runs </td><td>  -  </td></tr>
        <tr><td> 401 </td><td> Unauthorized </td><td>  -  </td></tr>
        <tr><td> 404 </td><td> Resource Not Found </td><td>  -  </td></tr>
        <tr><td> 0 </td><td> Internal Server Error </td><td>  -  </td></tr>
     </table>
     */
    public ApiResponse<ActionRunList> listRepositoryRunsWithHttpInfo(String repository, String after, Integer amount, String branch, String commit) throws ApiException {
        okhttp3.Call localVarCall = listRepositoryRunsValidateBeforeCall(repository, after, amount, branch, commit, null);
        Type localVarReturnType = new TypeToken<ActionRunList>(){}.getType();
        return localVarApiClient.execute(localVarCall, localVarReturnType);
    }

    /**
     * list runs (asynchronously)
     * 
     * @param repository  (required)
     * @param after return items after this value (optional)
     * @param amount how many items to return (optional, default to 100)
     * @param branch  (optional)
     * @param commit  (optional)
     * @param _callback The callback to be executed when the API call finishes
     * @return The request call
     * @throws ApiException If fail to process the API call, e.g. serializing the request body object
     * @http.response.details
     <table summary="Response Details" border="1">
        <tr><td> Status Code </td><td> Description </td><td> Response Headers </td></tr>
        <tr><td> 200 </td><td> list action runs </td><td>  -  </td></tr>
        <tr><td> 401 </td><td> Unauthorized </td><td>  -  </td></tr>
        <tr><td> 404 </td><td> Resource Not Found </td><td>  -  </td></tr>
        <tr><td> 0 </td><td> Internal Server Error </td><td>  -  </td></tr>
     </table>
     */
    public okhttp3.Call listRepositoryRunsAsync(String repository, String after, Integer amount, String branch, String commit, final ApiCallback<ActionRunList> _callback) throws ApiException {

        okhttp3.Call localVarCall = listRepositoryRunsValidateBeforeCall(repository, after, amount, branch, commit, _callback);
        Type localVarReturnType = new TypeToken<ActionRunList>(){}.getType();
        localVarApiClient.executeAsync(localVarCall, localVarReturnType, _callback);
        return localVarCall;
    }
    /**
     * Build call for listRunHooks
     * @param repository  (required)
     * @param runId  (required)
     * @param after return items after this value (optional)
     * @param amount how many items to return (optional, default to 100)
     * @param _callback Callback for upload/download progress
     * @return Call to execute
     * @throws ApiException If fail to serialize the request body object
     * @http.response.details
     <table summary="Response Details" border="1">
        <tr><td> Status Code </td><td> Description </td><td> Response Headers </td></tr>
        <tr><td> 200 </td><td> list specific run hooks </td><td>  -  </td></tr>
        <tr><td> 401 </td><td> Unauthorized </td><td>  -  </td></tr>
        <tr><td> 404 </td><td> Resource Not Found </td><td>  -  </td></tr>
        <tr><td> 0 </td><td> Internal Server Error </td><td>  -  </td></tr>
     </table>
     */
    public okhttp3.Call listRunHooksCall(String repository, String runId, String after, Integer amount, final ApiCallback _callback) throws ApiException {
        Object localVarPostBody = null;

        // create path and map variables
        String localVarPath = "/repositories/{repository}/actions/runs/{run_id}/hooks"
            .replaceAll("\\{" + "repository" + "\\}", localVarApiClient.escapeString(repository.toString()))
            .replaceAll("\\{" + "run_id" + "\\}", localVarApiClient.escapeString(runId.toString()));

        List<Pair> localVarQueryParams = new ArrayList<Pair>();
        List<Pair> localVarCollectionQueryParams = new ArrayList<Pair>();
        Map<String, String> localVarHeaderParams = new HashMap<String, String>();
        Map<String, String> localVarCookieParams = new HashMap<String, String>();
        Map<String, Object> localVarFormParams = new HashMap<String, Object>();

        if (after != null) {
            localVarQueryParams.addAll(localVarApiClient.parameterToPair("after", after));
        }

        if (amount != null) {
            localVarQueryParams.addAll(localVarApiClient.parameterToPair("amount", amount));
        }

        final String[] localVarAccepts = {
            "application/json"
        };
        final String localVarAccept = localVarApiClient.selectHeaderAccept(localVarAccepts);
        if (localVarAccept != null) {
            localVarHeaderParams.put("Accept", localVarAccept);
        }

        final String[] localVarContentTypes = {
            
        };
        final String localVarContentType = localVarApiClient.selectHeaderContentType(localVarContentTypes);
        localVarHeaderParams.put("Content-Type", localVarContentType);

        String[] localVarAuthNames = new String[] { "basic_auth", "cookie_auth", "jwt_token", "oidc_auth" };
        return localVarApiClient.buildCall(localVarPath, "GET", localVarQueryParams, localVarCollectionQueryParams, localVarPostBody, localVarHeaderParams, localVarCookieParams, localVarFormParams, localVarAuthNames, _callback);
    }

    @SuppressWarnings("rawtypes")
    private okhttp3.Call listRunHooksValidateBeforeCall(String repository, String runId, String after, Integer amount, final ApiCallback _callback) throws ApiException {
        
        // verify the required parameter 'repository' is set
        if (repository == null) {
            throw new ApiException("Missing the required parameter 'repository' when calling listRunHooks(Async)");
        }
        
        // verify the required parameter 'runId' is set
        if (runId == null) {
            throw new ApiException("Missing the required parameter 'runId' when calling listRunHooks(Async)");
        }
        

        okhttp3.Call localVarCall = listRunHooksCall(repository, runId, after, amount, _callback);
        return localVarCall;

    }

    /**
     * list run hooks
     * 
     * @param repository  (required)
     * @param runId  (required)
     * @param after return items after this value (optional)
     * @param amount how many items to return (optional, default to 100)
     * @return HookRunList
     * @throws ApiException If fail to call the API, e.g. server error or cannot deserialize the response body
     * @http.response.details
     <table summary="Response Details" border="1">
        <tr><td> Status Code </td><td> Description </td><td> Response Headers </td></tr>
        <tr><td> 200 </td><td> list specific run hooks </td><td>  -  </td></tr>
        <tr><td> 401 </td><td> Unauthorized </td><td>  -  </td></tr>
        <tr><td> 404 </td><td> Resource Not Found </td><td>  -  </td></tr>
        <tr><td> 0 </td><td> Internal Server Error </td><td>  -  </td></tr>
     </table>
     */
    public HookRunList listRunHooks(String repository, String runId, String after, Integer amount) throws ApiException {
        ApiResponse<HookRunList> localVarResp = listRunHooksWithHttpInfo(repository, runId, after, amount);
        return localVarResp.getData();
    }

    /**
     * list run hooks
     * 
     * @param repository  (required)
     * @param runId  (required)
     * @param after return items after this value (optional)
     * @param amount how many items to return (optional, default to 100)
     * @return ApiResponse&lt;HookRunList&gt;
     * @throws ApiException If fail to call the API, e.g. server error or cannot deserialize the response body
     * @http.response.details
     <table summary="Response Details" border="1">
        <tr><td> Status Code </td><td> Description </td><td> Response Headers </td></tr>
        <tr><td> 200 </td><td> list specific run hooks </td><td>  -  </td></tr>
        <tr><td> 401 </td><td> Unauthorized </td><td>  -  </td></tr>
        <tr><td> 404 </td><td> Resource Not Found </td><td>  -  </td></tr>
        <tr><td> 0 </td><td> Internal Server Error </td><td>  -  </td></tr>
     </table>
     */
    public ApiResponse<HookRunList> listRunHooksWithHttpInfo(String repository, String runId, String after, Integer amount) throws ApiException {
        okhttp3.Call localVarCall = listRunHooksValidateBeforeCall(repository, runId, after, amount, null);
        Type localVarReturnType = new TypeToken<HookRunList>(){}.getType();
        return localVarApiClient.execute(localVarCall, localVarReturnType);
    }

    /**
     * list run hooks (asynchronously)
     * 
     * @param repository  (required)
     * @param runId  (required)
     * @param after return items after this value (optional)
     * @param amount how many items to return (optional, default to 100)
     * @param _callback The callback to be executed when the API call finishes
     * @return The request call
     * @throws ApiException If fail to process the API call, e.g. serializing the request body object
     * @http.response.details
     <table summary="Response Details" border="1">
        <tr><td> Status Code </td><td> Description </td><td> Response Headers </td></tr>
        <tr><td> 200 </td><td> list specific run hooks </td><td>  -  </td></tr>
        <tr><td> 401 </td><td> Unauthorized </td><td>  -  </td></tr>
        <tr><td> 404 </td><td> Resource Not Found </td><td>  -  </td></tr>
        <tr><td> 0 </td><td> Internal Server Error </td><td>  -  </td></tr>
     </table>
     */
    public okhttp3.Call listRunHooksAsync(String repository, String runId, String after, Integer amount, final ApiCallback<HookRunList> _callback) throws ApiException {

        okhttp3.Call localVarCall = listRunHooksValidateBeforeCall(repository, runId, after, amount, _callback);
        Type localVarReturnType = new TypeToken<HookRunList>(){}.getType();
        localVarApiClient.executeAsync(localVarCall, localVarReturnType, _callback);
        return localVarCall;
    }
}
