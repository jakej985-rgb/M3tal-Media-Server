import sys
import os
from pathlib import Path

# --- M3TAL Path Bootstrap ---
p = Path(__file__).resolve()
root = None
for parent in [p] + list(p.parents):
    if (parent / ".env").exists() and (parent / "docker").exists():
        root = parent
        if str(parent / "dashboard") not in sys.path:
             sys.path.append(str(parent / "dashboard"))
        break

if not root:
    print("❌ Error: Could not find M3TAL project root.")
    sys.exit(1)

# Import from auth.py in dashboard
try:
    from auth import prompt_password, hash_password, load_users, save_users, resolve_users_path
except ImportError:
    print("❌ Error: Could not import auth module from dashboard.")
    sys.exit(1)

def main():
    print("\n" + "="*50)
    print(" 🔐 M3TAL DASHBOARD IDENTITY MANAGER")
    print("="*50)
    
    try:
        username = input("\n[?] Target Username [admin]: ").strip() or "admin"
        password = prompt_password(f"[?] Enter password for '{username}'")
        
        role = input("[?] Assign Role (admin/viewer) [admin]: ").strip().lower() or "admin"
        if role not in ["admin", "viewer"]:
            print(f"[!] Invalid role '{role}', defaulting to 'admin'")
            role = "admin"

        users_path = resolve_users_path()
        users = load_users(users_path=users_path)
        
        # Remove existing user to avoid duplicates
        original_count = len(users)
        users = [u for u in users if u["username"] != username]
        
        # Add the new/updated record
        users.append({
            "username": username,
            "token_hash": hash_password(password),
            "role": role
        })
        
        # Sort for consistency
        users.sort(key=lambda u: u["username"])
        
        save_users(users, users_path=users_path)
        
        action = "Updated" if len(users) == original_count else "Added"
        print(f"\n✅ {action} user '{username}' ({role})")
        print(f"📄 Registry: {users_path}")
        print("="*50 + "\n")
        
    except KeyboardInterrupt:
        print("\n\n[!] Operation cancelled.")
        sys.exit(1)
    except Exception as e:
        print(f"\n❌ Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
