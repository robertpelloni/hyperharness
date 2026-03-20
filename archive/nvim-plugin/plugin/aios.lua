if vim.g.loaded_borg then
    return
end
vim.g.loaded_borg = true

vim.api.nvim_create_user_command("BorgSetup", function()
    require("borg").setup()
end, { desc = "Setup Borg plugin" })
