# Educational LSP 

Exploring the inner workings of the language server protocol. 

## Getting Started

Build the project using ```make build```. Logs are output to the same directory. 

### Register LSP with Nvim 
I've setup my local environment to use this LSP for *.md files. This can be adjusted / changed @ <TODO> 


```lua
local client = vim.lsp.start_client {
    name = "educationallsp",
    cmd = { "<path>/<to>/<binary>/test-lsp" }
}

if not client then
    vim.notify "educationallsp did not start"
    return
end

vim.api.nvim_create_autocmd("FileType", {
    pattern = "markdown",
    callback = function()
        vim.lsp.buf_attach_client(0, client)
    end

})



```

## References
 
- TJ DeVries' LSP example project: https://www.youtube.com/watch?v=YsdlcQoHqPY&t=1715s
