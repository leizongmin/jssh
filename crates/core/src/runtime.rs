use anyhow::anyhow;
use rquickjs::{
    async_with,
    loader::{
        BuiltinLoader, BuiltinResolver, FileResolver, ModuleLoader, NativeLoader, ScriptLoader,
    },
    AsyncContext, AsyncRuntime, CatchResultExt, Ctx, Error, Function, Module, Promise, Result,
};

use crate::bundle::{NativeModule, SCRIPT_MODULE};

struct RuntimeCore {
    rt: AsyncRuntime,
    ctx: AsyncContext,
}

impl RuntimeCore {
    pub async fn new() -> Self {
        let rt = AsyncRuntime::new().unwrap();
        let ctx = AsyncContext::full(&rt).await.unwrap();

        let resolver = (
            BuiltinResolver::default()
                .with_module("bundle/script_module")
                .with_module("bundle/native_module"),
            FileResolver::default()
                .with_path("./")
                .with_path("../../target/debug")
                .with_native(),
        );
        let loader = (
            BuiltinLoader::default().with_module("bundle/script_module", SCRIPT_MODULE),
            ModuleLoader::default().with_module("bundle/native_module", NativeModule),
            ScriptLoader::default(),
            NativeLoader::default(),
        );
        rt.set_loader(resolver, loader);

        RuntimeCore { rt, ctx }
    }

    pub async fn evaluate_module(&self, name: &str, script: &str) -> Result<()> {
        let ret = async_with!(self.ctx => |ctx| {
            let p = Module::evaluate(ctx.clone(), name, script).catch(&ctx).map_err(|e| {
                println!("failed to evaluate user script: {:?}", e);
                anyhow!("failed to evaluate user script: {:?}", e)
           }).unwrap();
            let r = p.into_future::<()>();
        });
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use crate::runtime::*;

    #[tokio::test]
    async fn it_works() {
        let rt = RuntimeCore::new().await;
        rt.evaluate_module("repl", "1+1+2; throw new Error('hello')")
            .await.unwrap();
    }
}
