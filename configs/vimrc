""execute pathogen#infect()

syntax on
filetype plugin indent on

autocmd FileType yaml setlocal ai et sw=2 ts=2 sts=2 nu cuc

set foldlevelstart=20

let g:ale_echo_msg_format = '[%linter%] %s [%severity%]'
let g:ale_sign_error = '✘'
let g:ale_sign_warning = '⚠'
let g:ale_lint_on_text_changed = 'never'


if v:lang =~ "utf8$" || v:lang =~ "UTF-8$"
   set fileencodings=ucs-bom,utf-8,latin1
endif
"
set nocompatible	" Use Vim defaults (much better!)
set bs=indent,eol,start		" allow backspacing over everything in insert mode
""set ai			" always set autoindenting on
""set backup		" keep a backup file
set viminfo='20,\"50	" read/write a .viminfo file, don't store more
"			" than 50 lines of registers
set history=50		" keep 50 lines of command line history
set ruler		" show the cursor position all the time
"
"" Only do this part when compiled with support for autocommands
if has("autocmd")
  augroup redhat
  autocmd!
"  " In text files, always limit the width of text to 78 characters
"  " autocmd BufRead *.txt set tw=78
"  " When editing a file, always jump to the last cursor position
  autocmd BufReadPost *
  \ if line("'\"") > 0 && line ("'\"") <= line("$") |
  \   exe "normal! g'\"" |
  \ endif
"  " don't write swapfile on most commonly used directories for NFS mounts or USB sticks
  autocmd BufNewFile,BufReadPre /media/*,/run/media/*,/mnt/* set directory=~/tmp,/var/tmp,/tmp
"  " start with spec file template
  autocmd BufNewFile *.spec 0r /usr/share/vim/vimfiles/template.spec
  augroup END
endif
"
if has("cscope") && filereadable("/usr/bin/cscope")
   set csprg=/usr/bin/cscope
   set csto=0
   set cst
   set nocsverb
"   " add any database in current directory
   if filereadable("cscope.out")
	   cs add $PWD/cscope.out
   " else add database pointed to by environment
   elseif $CSCOPE_DB != ""
      cs add $CSCOPE_DB
   endif
   set csverb
endif
"
"" Switch syntax highlighting on, when the terminal has colors
"" Also switch on highlighting the last used search pattern.
if &t_Co > 2 || has("gui_running")
  syntax on
  set hlsearch
  endif
"
filetype plugin on
"
if &term=="xterm"
     set t_Co=8
     set t_Sb=[4%dm
     set t_Sf=[3%dm
endif
"
"" Don't wake up system with blinking cursor:
"
"" http://www.linuxpowertop.org/known.php
let &guicursor = &guicursor . ",a:blinkon0"
"
"############################################USER DEFINED OPTIONS, MAPS,FUNCTIONS, ETC##############################################
"
if &t_Co > 2 || has("gui_running")
  syntax on
  set hlsearch
  endif
"source /etc/vimrc
"
"}}}

"COLORSCHEME AND OPTIONS-----{{{

colorscheme delek
set number
set autoindent
filetype plugin indent on
set splitright
set tildeop
set splitbelow
set showcmd
set smartcase
set ignorecase
"set smartindent
set linebreak
set showbreak=---->
set breakindent
set virtualedit=block,onemore
set laststatus=2
set foldlevelstart=0
set tabstop=4
set shiftwidth=4
set expandtab
set showtabline=1
highlight cursorline cterm=underline
highlight Folded Ctermbg=darkred ctermfg=white
set foldmethod=marker
"set cursorline
"let first_line = getline(1)
"if first_line == 'notes'
"	setfiletype notes
"	"echom getline(1)
"endif
"}}}

"GLOBAL VARIABLE DEFINITIONS-----{{{

let mapleader = " "
let localmapleader = "\\"

"}}}

"NORMAL MODE GLOBAL NON-RECURSIVE MAPPINGS-----{{{

nnoremap - ddp
nnoremap <expr> _ line(".")==1 ? "" : "ddk<s-p>"
nnoremap <enter> i<enter><esc>
nnoremap <leader>j 10j
nnoremap <leader>k 10k
nnoremap <c-j> <pagedown>

nnoremap <c-k> <pageup>
nnoremap <c-u> viw<s-u>
nnoremap <leader>sv :source $MYVIMRC<cr>
nnoremap <leader>" viw<esc>a"<esc>hbi"<esc>lel
nnoremap <leader>' viw<esc>a'<esc>hbi'<esc>lel
nnoremap <leader>h 0
nnoremap <leader>l $
nnoremap <up> <nop>
nnoremap <down> <nop>
nnoremap <left> <nop>
nnoremap <right> <nop>
nnoremap <leader>eb :call EditBash()<cr>
nnoremap dw lbdw
nnoremap dW dt<space>
nnoremap <leader>ev :call EditVim()<cr>
nnoremap <leader>c @='0i<c-v><esc>l'<cr>
nnoremap noh :noh<cr>
nnoremap <BS> i<BS><esc>l
nnoremap <leader>w :write<cr>
nnoremap <leader>q :q!<cr>
nnoremap <leader>a :qa!<cr>
nnoremap dl 0d$
nnoremap q<leader> q:
nnoremap / /\v
nnoremap \ :call WindowSearch("
nnoremap t gt
nnoremap T gT
"nnoremap sf i<space>"---{{{<esc>bli
"nnoremap ef i<space>"}}}<esc>0
nnoremap sf i"---{{{<esc>bli
nnoremap ef i"}}}<esc>0
nnoremap bf i#---{{{<esc>bli
nnoremap cf i#}}}<esc>0
nnoremap <tab> i<tab><esc>
"nnoremap <c-q> <c-u>normal! q:

"}}}

"VISUAL MODE GLOBAL NON-RECURSIVE MAPPINGS-----{{{

vnoremap <leader>c 0<s-i>#<esc>
vnoremap fd <esc>
vnoremap <esc> <nop>
vnoremap <s-q> <esc>`<i"<esc>`>i"<esc>

"}}}

"INSERT MODE GLOBAL NON-RECURSIVE MAPPINGS AND ABBREVIATIONS-----{{{

inoremap jk <esc>
inoremap <c-u> <esc>viw<s-u>i
inoremap <c-z> <esc><c-z>
inoremap <c-Z> <esc><c-z>
execute "set <M-j>=\ej"
execute "set <M-k>=\ek"
execute "set <M-l>=\el"
execute "set <M-h>=\eh"
"execute "set <M-0=\e0"
inoremap <M-j> <esc>ji
inoremap <M-k> <esc>ki
inoremap <M-h> <esc>i
inoremap <M-l> <esc>lli
inoremap <M-0> <esc>0i
"inoremap <esc> <nop>

"}}}

"EX MODE GLOBAL NON-RECURSIVE MAPPINGS---{{{

"command te tabedit
"cmap te tabe
cnoreabbrev <expr> te ((getcmdtype() is# ':' && getcmdline() is# 'te')?('tabe'):('te'))
"
"
"
"}}}

"OPERATOR PENDING MAPPINGS-----{{{

onoremap inp :<c-u>normal! f(vi(<cr>
onoremap ilp :<c-u>normal! F)vi(<cr>
onoremap anp :<c-u>normal! f(va(<cr>
onoremap alp :<c-u>normal! F)va(<cr>
onoremap cp :<c-u>normal! F(vi(<cr>

onoremap inb :<c-u>normal! f[vi[<cr>
onoremap ilb :<c-u>normal! F]vi[<cr>
onoremap anb :<c-u>normal! f[va[<cr>
onoremap alb :<c-u>normal! F]va[<cr>
onoremap cb :<c-u>normal! F[vi[<cr>

onoremap inc :<c-u>normal! f{vi{<cr>
onoremap ilc :<c-u>normal! F}vi{<cr>
onoremap anc :<c-u>normal! f{va{<cr>
onoremap alc :<c-u>normal! F}va{<cr>
onoremap cc :<c-u>normal! F{vi{<cr>

"}}}

"AUGROUP DEFINITIONS-----{{{

augroup shell
	au!
	autocmd FileType sh nnoremap <buffer> <localleader>c I<esc>
	autocmd FileType sh iabbrev <buffer> 'fun' function 
augroup END

augroup python 
	au!
	autocmd filetype python nnoremap <buffer> <localleader>c <s-i><esc>
augroup END

augroup c_commands
 	au!
 	autocmd FileType c nnoremap <buffer> <localleader>c <s-i>//<esc>
augroup END

augroup filetype_vim
	au!
	autocmd FileType vim setlocal foldmethod=marker
augroup END

augroup filetype_notes
	au!
	autocmd FileType notes setlocal foldmethod=marker
    autocmd Filetype notes setlocal syntax=vim
    autocmd Filetype notes setlocal smartindent
augroup END

augroup filetype_qf
	au!
	autocmd Filetype qf execute 'nnoremap <buffer> <enter> <enter>'
augroup END

augroup filetype_journal
	au!
	autocmd Filetype journal silent execute '!chmod 777 /tmp/journal'
	autocmd Filetype journal execute 'write! /tmp/journal'
augroup END

"}}}

"AUTOCMD GROUP SOURCING-----{{{

autocmd FileType sh doautocmd BufReadPost shell
autocmd FileType python doautocmd BufReadPost python
autocmd Filetype c doautocmd BufReadPost c_commands
autocmd InsertEnter,InsertLeave * set cul!
autocmd FileType help wincmd L
autocmd BufNewFile,BufRead * if getline(1) =~ '^\/\* notes' | setlocal filetype=notes | endif
autocmd BufNewFile,BufRead * if getline(1) =~ '^Notes' | setlocal filetype=notes | endif
autocmd BufNewFile,BufRead,StdinReadPost * if getline(1) =~ '^-- Logs begin at' | setlocal filetype=journal | endif
autocmd Filetype notes setlocal foldmethod=marker
"autocmd Bufread getline(1)

"}}}

"FUNCTIONS-----{{{

function EditVim()
	set nosplitright
	vsp $MYVIMRC
	set splitright
endfunction

function EditBash()
	set nosplitright
	vsp $profile
	set splitright
endfunction

function SetStatusLine()
	set statusline=%F
	set statusline+=\ \ \ Filetype:\ %y
	set statusline+=\ \ \ Line:\ %l/%L
    set statusline+=\ \ \ Column:\ %v
    set statusline+=\ \ \ Character\ Code:\ %b
endfunction
call SetStatusLine()

"function SetTabLine()
"    set tabline=%N\j
"endfunction
"call SetTabLine()

function WindowSearch(pattern)
	execute 'vimgrep' a:pattern '% | cwindow | wincmd L'
endfunction


""Command Output To Left Window---{{{
"
"command! -complete=shellcmd -nargs=+ Shell call s:RunShellCommand(<q-args>)
"
"function! s:RunShellCommand(cmdline)
"	"echo a:cmdline
"	let expanded_cmdline = a:cmdline
"	for part in split(a:cmdline, ' ')
"		if part[0] =~ '\v[%#<]'
"			let expanded_part = fnameescape(expand(part))
"			let expanded_cmdline = substitute(expanded_cmdline, part, expanded_part, '')
"			endif
"	endfor
"	botright new
"	setlocal bufhidden=wipe nobuflisted noswapfile nowrap filetype=command
"	call setline(1, 'You entered:    ' . a:cmdline)
"	call setline(2, 'Expanded Form:  ' .expanded_cmdline)
"	call setline(3,substitute(getline(2),'.','=','g'))
"	execute '$read !'. expanded_cmdline
"	"setlocal nomodifiable
"	wincmd L 
"	1
"endfunction
"
""}}}
"
""}}}
"
""OLD/OBSOLETE/TEST MAPPINGS/ABBREVIATIONS-----{{{
"	"test_of_fold_level_2----{{{
"	"
"	"
"	"
"	"
"	"
"	"}}}
"" nnoremap - 0d$jp
""iabbrev kyle loser
""nnoremap <expr> <leader>__ line(".")==1 ? "" : "ddk<s-p>"
"" echo (>^.^<)
""inoremap <up> ""
""inoremap <down> <nop>
""inoremap <left> <nop>
""nnoremap <leader>ev :vsp $MYVIMRC<cr>
""nnoremap <leader>- ddp
""nnoremap _ dd<s-p>
""nnoremap <leader>_ i<esc> v:count1 
""inoremap } <esc>lli
""inoremap { {}<esc>i
""inoremap <right> <nop>
""vnoremap <s-q> d`<,`>
" "return system.string(string1, string2); if (x <=5); done
" "
" "}}}
